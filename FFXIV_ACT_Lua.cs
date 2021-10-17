/*
This file is part of FFXIV ACT Lua.

FFXIV ACT Lua is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FFXIV ACT Lua is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FFXIV ACT Lua.  If not, see <https://www.gnu.org/licenses/>.
*/

using System;
using System.IO;
using System.Text;
using System.Collections.Generic;
using System.Windows.Forms;
using System.Reflection;
using System.Net;
using System.Net.Sockets;
using System.Threading;  
using System.Threading.Tasks;
using Advanced_Combat_Tracker;

[assembly: AssemblyTitle("FFXIV ACT Lua")]
[assembly: AssemblyDescription("Bup!")]
[assembly: AssemblyCompany("Nathan Ogden")]
[assembly: AssemblyVersion("0.01")]

namespace ACT_Plugin
{
    public class FFActLua : UserControl, IActPluginV1
    {

        const int VERSION_NUMBER = 1;

        const UInt16 DAEMON_PORT = 31593;                       // Port to send to daemon on.
        
        const byte DATA_TYPE_SESSION = 1;                       // Data type, session data
        const byte DATA_TYPE_ENCOUNTER = 2;                     // Data type, encounter data
        const byte DATA_TYPE_COMBATANT = 3;                     // Data type, combatant data
        const byte DATA_TYPE_LOG_LINE = 5;                      // Data type, log line
        const byte DATA_TYPE_FLAG = 99;                         // Data type, flag

        const byte DATA_TYPE_SCRIPTS_AVAILABLE = 201;           // Data type, lua scripts available
        const byte DATA_TYPE_SCRIPTS_ENABLED = 202;             // Data type, lua scripts enabled
        const byte DATA_TYPE_ACT_SAY = 203;                     // Data type, speak with TTS
        const byte DATA_TYPE_ACT_END = 204;                     // Data type, flag to end encounter

        const long TTS_TIMEOUT = 3000;                          // Time in miliseconds to timeout TTS
        
        private Label lblStatus;                                // The status label that appears in ACT's Plugin tab
        private UdpClient udpClient;                            // UDP client used to send data
        private UdpClient udpListener;                          // UDP listener used to recv data
        private IPEndPoint udpEndpoint;                         // UDP address
        Thread listenThread;                                    // Thread for listening for incoming data
        private long lastTTSTime = 0;                           // Last time TTS was timed out
        private string[] availableScripts;
        private string[] enabledScripts;

        public FFActLua()
        {
        }

        public void InitPlugin(TabPage pluginScreenSpace, Label pluginStatusText)
        {
            // status label
            lblStatus = pluginStatusText; // Hand the status label's reference to our local var
            // update plugin status text
            lblStatus.Text = "Plugin version " + ((float) VERSION_NUMBER / 100.0).ToString("n2") + " started.";
            // init udp client
            udpConnect();
            // hook events
            ActGlobals.oFormActMain.OnCombatStart += new CombatToggleEventDelegate(oFormActMain_OnCombatStart);
            ActGlobals.oFormActMain.OnCombatEnd += new CombatToggleEventDelegate(oFormActMain_OnCombatEnd);
            ActGlobals.oFormActMain.AfterCombatAction += new CombatActionDelegate(oFormActMain_AfterCombatAction);
            ActGlobals.oFormActMain.OnLogLineRead += new LogLineEventDelegate(oFormActMain_OnLogLineRead);
            // form stuff
			//pluginScreenSpace.Controls.Add(this);	// Add this UserControl to the tab ACT provides
			//this.Dock = DockStyle.Fill;	// Expand the UserControl to fill the tab's client space
        }

        public void DeInitPlugin()
        {
            // deinit event hooks
            ActGlobals.oFormActMain.OnCombatStart -= oFormActMain_OnCombatStart;
            ActGlobals.oFormActMain.OnCombatEnd -= oFormActMain_OnCombatEnd;
            ActGlobals.oFormActMain.AfterCombatAction -= oFormActMain_AfterCombatAction;
            ActGlobals.oFormActMain.OnLogLineRead  -= oFormActMain_OnLogLineRead;
            //this.buttonSave.Click -= buttonSave_OnClick;
            // close udp client
            udpClient.Close();
            if (listenThread != null) {
                listenThread.Abort();
            }
            udpListener.Close();
            // update plugin status text
            lblStatus.Text = "Plugin Exited";
        }

        void oFormActMain_OnCombatStart(bool isImport, CombatToggleEventArgs actionInfo)
        {
            sendEncounterData(actionInfo.encounter);
        }

        void oFormActMain_OnCombatEnd(bool isImport, CombatToggleEventArgs actionInfo)
        {
            sendEncounterData(actionInfo.encounter);
        }

        void oFormActMain_AfterCombatAction(bool isImport, CombatActionEventArgs actionInfo)
        {
            sendCombatData(actionInfo);
        }

        void oFormActMain_OnLogLineRead(bool isImport, LogLineEventArgs logInfo)
        {
            sendLogLine(logInfo);
        }

        void udpConnect()
        {
            // connect to daemon
            if (udpClient != null) {
                udpClient.Close();
            }
            udpEndpoint = new IPEndPoint(IPAddress.Parse("127.0.0.1"), DAEMON_PORT);
            udpClient = new UdpClient();
            try {
                udpClient.Connect(udpEndpoint);
            } catch (System.Net.Sockets.SocketException e) {
                lblStatus.Text = "ERROR: " + e.Message;
            }
            // listen from daemon
            if (udpListener != null) {
                udpListener.Close();
            }
            udpListener = new UdpClient();
            listenThread = new Thread(new ThreadStart(udpListen));
            listenThread.Start();
        }

        void sendUdp(ref List<Byte> sendData)
        {
            Byte[] sendBytes = sendData.ToArray();
            udpClient.Send(sendBytes, sendBytes.Length);          
        }

        void udpListen()
        {
            while(true)
            {
                try {
                    byte[] bytes = udpClient.Receive(ref udpEndpoint);
                    switch (bytes[0]) {
                        case DATA_TYPE_ACT_SAY: {
                            string text = System.Text.Encoding.UTF8.GetString(bytes, 1, bytes.Length-1);
                            ttsSay(text);
                            break;
                        }
                        case DATA_TYPE_ACT_END: {
                            ActGlobals.oFormActMain.EndCombat(true);
                            break;
                        }
                        case DATA_TYPE_SCRIPTS_AVAILABLE: {
                            string valStr = System.Text.Encoding.UTF8.GetString(bytes, 1, bytes.Length-1);
                            availableScripts = valStr.Split(',');
                            break;
                        }
                        case DATA_TYPE_SCRIPTS_ENABLED: {
                            string valStr = System.Text.Encoding.UTF8.GetString(bytes, 1, bytes.Length-1);
                            enabledScripts = valStr.Split(',');
                            break;
                        }
                    }

                } catch (SocketException e) {
                    lblStatus.Text = "ERROR: " + e.Message;
                }
            }
            //udpListener.Close();
        }

        void prepareDateTime(ref List<Byte> sendData, DateTime value) 
        {
            string dateTimeString = value.ToString("o");
            prepareString(ref sendData, dateTimeString);
        }

        void prepareUint16(ref List<Byte> sendData, UInt16 value)
        {
            Byte[] valueBytes = BitConverter.GetBytes((UInt16)value);
            if (BitConverter.IsLittleEndian) {
                 Array.Reverse(valueBytes);
            }
            sendData.AddRange(valueBytes);   
        }

        void prepareInt32(ref List<Byte> sendData, Int32 value)
        {
            Byte[] valueBytes = BitConverter.GetBytes((Int32)value);
            if (BitConverter.IsLittleEndian) {
                 Array.Reverse(valueBytes);
            }
            sendData.AddRange(valueBytes);
        }

        void prepareString(ref List<Byte> sendData, string value)
        {
            Byte[] valueBytes = Encoding.UTF8.GetBytes(value);
            prepareUint16(ref sendData, (UInt16) valueBytes.Length);
            if (valueBytes.Length > 0) {
                sendData.AddRange(valueBytes);
            }
        }

        void sendCombatData(CombatActionEventArgs actionInfo)
        {
            if (actionInfo.cancelAction) {
                return;
            }
            // get encounter
            EncounterData encounter = actionInfo.combatAction.ParentEncounter;
            // send encounter data
            sendEncounterData(encounter);
            // send combatant data
            foreach (CombatantData cd in encounter.GetAllies()) {
                if (cd.Name == actionInfo.attacker) {
                    // get actor id (stored as tag in actionInfo)
                    Int32 actorId = int.Parse(
                        (string) actionInfo.tags["ActorID"],
                        System.Globalization.NumberStyles.HexNumber
                    );
                    // send combatant data with actor id
                    sendEncounterCombatantData(cd, actorId);
                    break;
                }
            }
        }

        void sendEncounterData(EncounterData ed)
        {
            // build send data
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_ENCOUNTER);                             // declare data type
            prepareInt32(ref sendData, ed.StartTime.GetHashCode());        // encounter id (start time hash code)
            prepareDateTime(ref sendData, ed.StartTime);                   // start time of encounter
            prepareDateTime(ref sendData, ed.EndTime);                     // end time of encounter
            prepareString(ref sendData, ed.ZoneName);                      // zone name
            prepareInt32(ref sendData, (Int32) ed.Damage);                 // encounter damage
            sendData.Add((byte) (ed.Active ? 1 : 0));                      // is still active encounter
            sendData.Add((byte) ed.GetEncounterSuccessLevel());            // success level of encounter
            // send
            sendUdp(ref sendData);
        }

        void sendEncounterCombatantData(CombatantData cd, Int32 actorId)
        {           
            // build send data
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_COMBATANT);                             // declare data type
            prepareInt32(ref sendData, cd.EncStartTime.GetHashCode());     // encounter id
            prepareInt32(ref sendData, actorId);                           // actor id (ffxiv)
            prepareString(ref sendData, cd.Name);                          // combatant name
            prepareString(ref sendData, cd.GetColumnByName("Job"));        // combatant job (ffxiv)
            prepareInt32(ref sendData, (Int32) cd.Damage);                 // damage done
            prepareInt32(ref sendData, (Int32) cd.DamageTaken);            // damage taken
            prepareInt32(ref sendData, (Int32) cd.Healed);                 // damage healed
            prepareInt32(ref sendData, cd.Deaths);                         // number of deaths
            prepareInt32(ref sendData, cd.Hits);                           // number of attacks
            prepareInt32(ref sendData, cd.Heals);                          // number of heals performed
            prepareInt32(ref sendData, cd.Kills);                          // number of kills
            // send
            sendUdp(ref sendData);
        }

        void sendLogLine(LogLineEventArgs logInfo)
        {
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_LOG_LINE);
            // encounter id, if active
            Int32 encounterId = 0;
            if (logInfo.inCombat) {
                encounterId = ActGlobals.oFormActMain.ActiveZone.ActiveEncounter.StartTime.GetHashCode();
            }
            prepareInt32(ref sendData, encounterId);
            // time
            prepareDateTime(ref sendData, logInfo.detectedTime);
            // line
            prepareString(ref sendData, logInfo.logLine);
            // send
            sendUdp(ref sendData);
        }

        void ttsSay(string text)
        {
            long now = DateTime.Now.Ticks / TimeSpan.TicksPerMillisecond;
            if (lastTTSTime == 0 || now - lastTTSTime > TTS_TIMEOUT) {
                lastTTSTime = now;
                ActGlobals.oFormActMain.TTS(text);
            }
        }

    }
}
