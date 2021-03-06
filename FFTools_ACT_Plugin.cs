/*
This file is part of FFTools.

FFTools is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

FFTools is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with FFTools.  If not, see <https://www.gnu.org/licenses/>.
*/

using System;
using System.IO;
using System.Text;
using System.Collections.Generic;
using System.Windows.Forms;
using System.Reflection;
using System.Net;
using System.Net.Sockets;
using System.Diagnostics;
using System.Threading;  
using System.Threading.Tasks;
using Advanced_Combat_Tracker;

[assembly: AssemblyTitle("FFTools")]
[assembly: AssemblyDescription("Extends FFXIV parsing with Lua scripts that support TTS callouts, web UI, and more.")]
[assembly: AssemblyCompany("Chompy#3436")]
[assembly: AssemblyVersion("0.10")]

namespace ACT_Plugin
{
    public class FFTools : UserControl, IActPluginV1
    {

        const int VERSION_NUMBER = 10;

        const UInt16 DAEMON_PORT = 31593;                       // Port to send to daemon on.
        
        const byte DATA_TYPE_ENCOUNTER = 2;                     // Data type, encounter data
        const byte DATA_TYPE_COMBATANT = 3;                     // Data type, combatant data
        const byte DATA_TYPE_LOG_LINE = 5;                      // Data type, log line

        const byte DATA_TYPE_SCRIPT = 201;                      // Data type, information about an available lua script
        const byte DATA_TYPE_SCRIPT_ENABLE = 202;               // Data type, enable script
        const byte DATA_TYPE_SCRIPT_DISABLE = 203;              // Data type, disable script
        const byte DATA_TYPE_SCRIPT_RELOAD = 204;               // Data type, reload script
        const byte DATA_TYPE_SCRIPT_VERSION = 205;              // Data type, request/recieve script version
        const byte DATA_TYPE_SCRIPT_UPDATE = 206;               // Data type, request script update
        const byte DATA_TYPE_ACT_SAY = 210;                     // Data type, speak with TTS
        const byte DATA_TYPE_ACT_END = 211;                     // Data type, flag to end encounter
        const byte DATA_TYPE_ACT_UPDATE = 212;                  // Data type, flag that an update is ready

        const long TTS_TIMEOUT = 500;                           // Time in miliseconds to timeout TTS
        
        private Label lblStatus;                                // The status label that appears in ACT's Plugin tab
        private UdpClient udpClient;                            // UDP client used to send data
        private UdpClient udpListener;                          // UDP listener used to recv data
        private IPEndPoint udpEndpoint;                         // UDP address
        Thread listenThread;                                    // Thread for listening for incoming data
        private long lastTTSTime = 0;                           // Last time TTS was timed out
        private List<string[]> scriptData;                      // List of available Lua scripts
        private string lastScriptSelected;                      // Name of last script selected in list
        private Process scriptDaemon;                           // Instance of script daemon process
        private Thread scriptDaemonExitThread;                  // Thread to listen for script daemon exit

        private System.Windows.Forms.ListBox formScriptList;    // Form element containing list of available Lua scripts
        private System.Windows.Forms.TextBox formScriptInfo;    // Form element containing information about selected script
        private System.Windows.Forms.Button formScriptEnable;   // Form element button to enable/disable script
        private System.Windows.Forms.Button formScriptConfig;   // Form element button to open script config file in notepad
        private System.Windows.Forms.Button formScriptReload;   // Form element button to reload scripts
        private System.Windows.Forms.Button formWebOpen;        // Form element button to open web page for script
        private System.Windows.Forms.Button formUpdate;         // Form element button to check for script update
        private System.Windows.Forms.Button formScriptDir;      // Form element button to open script directory
        private System.Windows.Forms.Button formAppConfig;      // Form element button to open app config file in notepad

        
        public FFTools()
        {
            this.SuspendLayout();
            this.Dock = DockStyle.Fill;

            this.formScriptList = new System.Windows.Forms.ListBox();
            this.formScriptList.Name = "ScriptList";
            this.formScriptList.Location = new System.Drawing.Point(12, 12);
            this.formScriptList.MinimumSize = new System.Drawing.Size(245, 100);
            this.formScriptList.AutoSize = false;
            this.formScriptList.Dock = DockStyle.Bottom|DockStyle.Top;
            this.formScriptList.Font = new System.Drawing.Font("Consolas", 10);
            this.Controls.Add(this.formScriptList);

            this.formScriptInfo = new System.Windows.Forms.TextBox();
            this.formScriptInfo.Name = "ScriptInfo";
            this.formScriptInfo.Location = new System.Drawing.Point(252, 0);
            this.formScriptInfo.MinimumSize = new System.Drawing.Size(245, 200);
            this.formScriptInfo.AutoSize = false;
            this.formScriptInfo.ReadOnly = true;
            this.formScriptInfo.Multiline = true;
            this.Controls.Add(this.formScriptInfo);

            this.formScriptEnable = new System.Windows.Forms.Button();
            this.formScriptEnable.Name = "ScriptEnable";
            this.formScriptEnable.AutoSize = true;
            this.formScriptEnable.Location = new System.Drawing.Point(252, this.formScriptInfo.Height + 4);
            this.formScriptEnable.Size = new System.Drawing.Size(64, 24);
            this.formScriptEnable.Text = "Enable";
            this.formScriptEnable.Enabled = false;
            this.Controls.Add(this.formScriptEnable);

            this.formScriptConfig = new System.Windows.Forms.Button();
            this.formScriptConfig.Name = "ScriptConfig";
            this.formScriptConfig.AutoSize = true;
            this.formScriptConfig.Location = new System.Drawing.Point(318, this.formScriptInfo.Height + 4);
            this.formScriptConfig.Size = new System.Drawing.Size(82, 24);
            this.formScriptConfig.Text = "Edit Config";
            this.formScriptConfig.Enabled = false;
            this.Controls.Add(this.formScriptConfig);

            this.formWebOpen = new System.Windows.Forms.Button();
            this.formWebOpen.Name = "WebOpen";
            this.formWebOpen.AutoSize = true;
            this.formWebOpen.Location = new System.Drawing.Point(402, this.formScriptInfo.Height + 4);
            this.formWebOpen.Size = new System.Drawing.Size(82, 24);
            this.formWebOpen.Text = "Open Web View";
            this.formWebOpen.Enabled = false;
            this.Controls.Add(this.formWebOpen);

            this.formUpdate = new System.Windows.Forms.Button();
            this.formUpdate.Name = "ScriptUpdate";
            this.formUpdate.AutoSize = true;
            this.formUpdate.Location = new System.Drawing.Point(252, this.formScriptInfo.Height + 30);
            this.formUpdate.Size = new System.Drawing.Size(82, 24);
            this.formUpdate.Text = "Check For Updates";
            this.Controls.Add(this.formUpdate);

            this.formScriptReload = new System.Windows.Forms.Button();
            this.formScriptReload.Name = "ScriptReload";
            this.formScriptReload.AutoSize = true;
            this.formScriptReload.Location = new System.Drawing.Point(252, this.formScriptInfo.Height + 64);
            this.formScriptReload.Size = new System.Drawing.Size(82, 24);
            this.formScriptReload.Text = "Refresh List";
            this.formScriptReload.Enabled = true;
            this.Controls.Add(this.formScriptReload);

            this.formScriptDir = new System.Windows.Forms.Button();
            this.formScriptDir.Name = "ScriptDir";
            this.formScriptDir.AutoSize = true;
            this.formScriptDir.Location = new System.Drawing.Point(336, this.formScriptInfo.Height + 64);
            this.formScriptDir.Size = new System.Drawing.Size(114, 24);
            this.formScriptDir.Text = "Open Script Directory";
            this.Controls.Add(this.formScriptDir);

            this.formAppConfig = new System.Windows.Forms.Button();
            this.formAppConfig.Name = "AppConfig";
            this.formAppConfig.AutoSize = true;
            this.formAppConfig.Location = new System.Drawing.Point(252, this.formScriptInfo.Height + 90);
            this.formAppConfig.Size = new System.Drawing.Size(142, 24);
            this.formAppConfig.Text = "Edit Main Plugin Config";
            this.Controls.Add(this.formAppConfig);

            this.ResumeLayout(false);
            this.PerformLayout();
        }

        public void InitPlugin(TabPage pluginScreenSpace, Label pluginStatusText)
        {
            // start daemon
            this.scriptDaemon = new Process();
            this.scriptDaemon.StartInfo.FileName = this.getPluginDirectory() + "\\bin\\fftools_daemon.exe";
            this.scriptDaemon.StartInfo.CreateNoWindow = true;
            this.scriptDaemon.StartInfo.WindowStyle = System.Diagnostics.ProcessWindowStyle.Hidden;
            this.scriptDaemon.StartInfo.UseShellExecute = false;
            this.scriptDaemon.StartInfo.RedirectStandardError = true;
            this.scriptDaemon.Start();
            this.scriptDaemonExitThread = new Thread(new ThreadStart(this.listenDaemonExit));
            this.scriptDaemonExitThread.Start();
            Thread.Sleep(250);
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
            this.formScriptList.SelectedIndexChanged += ScriptList_SelectedIndexChanged;
            this.formScriptList.DoubleClick  += ScriptList_DoubleClick;
            this.formScriptEnable.Click += ScriptEnable_Click;
            this.formScriptReload.Click += ScriptReload_Click;
            this.formScriptConfig.Click += ScriptConfig_Click;
            this.formWebOpen.Click += ScriptWeb_Click;
            this.formScriptDir.Click += ScriptDir_Click;
            this.formAppConfig.Click += AppConfig_Click;
            this.formUpdate.Click += ScriptUpdate_Click;
            // /
			pluginScreenSpace.Controls.Add(this);	// Add this UserControl to the tab ACT provides
            // set tab title
            foreach (ActPluginData p in ActGlobals.oFormActMain.ActPlugins) {
                if (p.pluginObj == this) {
                    p.tpPluginSpace.Text = "FFTools";
                }
            }
        }

        public void DeInitPlugin()
        {
            // deinit event hooks
            ActGlobals.oFormActMain.OnCombatStart -= oFormActMain_OnCombatStart;
            ActGlobals.oFormActMain.OnCombatEnd -= oFormActMain_OnCombatEnd;
            ActGlobals.oFormActMain.AfterCombatAction -= oFormActMain_AfterCombatAction;
            ActGlobals.oFormActMain.OnLogLineRead  -= oFormActMain_OnLogLineRead;
            this.formScriptList.SelectedIndexChanged -= ScriptList_SelectedIndexChanged;
            this.formScriptEnable.Click -= ScriptEnable_Click;
            this.formScriptReload.Click -= ScriptReload_Click;
            this.formScriptConfig.Click -= ScriptConfig_Click;
            this.formWebOpen.Click -= ScriptWeb_Click;
            this.formScriptDir.Click -= ScriptDir_Click;
            this.formAppConfig.Click -= AppConfig_Click;
            this.formUpdate.Click -= ScriptUpdate_Click;
            // close udp client
            this.udpClient.Close();
            if (this.listenThread != null) {
                this.listenThread.Abort();
            }
            this.udpListener.Close();
            // abort script daemon exit thread
            if (this.scriptDaemonExitThread != null) {
                this.scriptDaemonExitThread.Abort();
            }
            // close process
            if (this.scriptDaemon != null) {
                this.scriptDaemon.Kill();
            }
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

        void ScriptList_SelectedIndexChanged(object sender, System.EventArgs e)
        {
            if (this.formScriptList.SelectedItem == null) {
                return;
            }
            string value = this.formScriptList.SelectedItem.ToString();
            int index = this.formScriptList.FindString(value);
            string titleSep = new String('-', 64);
            string desc = this.scriptData[index][3];
            if (this.scriptData[index][6] != "") {
                desc = "ERROR:\r\n" + this.scriptData[index][6];
            }
            this.formScriptInfo.Text = this.scriptData[index][2] + " (v" + this.scriptData[index][4] + ")\r\n" +
                titleSep + "\r\n" + desc;  
            this.formScriptEnable.Enabled = true;
            this.formScriptEnable.Text = "Disable";
            if (this.scriptData[index][1] == "") {
                this.formScriptEnable.Text = "Enable";
            }
            string configPath = this.getScriptConfigFile(this.scriptData[index][0]);
            this.formScriptConfig.Enabled = false;
            if (configPath != "") {
                this.formScriptConfig.Enabled = true;
            }
            this.formWebOpen.Enabled = true;
        }

        void ScriptList_DoubleClick(object sender, System.EventArgs e)
        {
            this.ScriptEnable_Click(sender, e);
        }

        void ScriptEnable_Click(object sender, System.EventArgs e)
        {
            this.formScriptList.Enabled = false;
            this.formScriptEnable.Enabled = false;
            string value = this.formScriptList.SelectedItem.ToString();
            int index = this.formScriptList.FindString(value);
            string name = this.scriptData[index][0];
            if (this.scriptData[index][1] == "") {
                this.resetScriptList();
                sendEnableScript(name);
                return;
            }
            this.resetScriptList();
            sendDisableScript(name);
        }

        void ScriptReload_Click(object sender, System.EventArgs e)
        {
            this.formScriptList.Enabled = false;
            this.formScriptEnable.Enabled = false;
            this.sendScriptRequest();
        }

        void ScriptConfig_Click(object sender, System.EventArgs e)
        {
            string value = this.formScriptList.SelectedItem.ToString();
            int index = this.formScriptList.FindString(value);
            string name = this.scriptData[index][0];
            string path = this.getScriptConfigFile(name);
            Process.Start("notepad.exe", path);
        }

        void ScriptWeb_Click(object sender, System.EventArgs e)
        {
            string value = this.formScriptList.SelectedItem.ToString();
            int index = this.formScriptList.FindString(value);
            string url = this.scriptData[index][5];
            Process.Start(url);
        }

        void ScriptDir_Click(object sender, System.EventArgs e)
        {
            string path = getPluginDirectory() + "\\scripts";
            Process.Start(path);
        }

        void ScriptUpdate_Click(object sender, System.EventArgs e)
        {
            string value = this.formScriptList.SelectedItem.ToString();
            int index = this.formScriptList.FindString(value);
            string name = this.scriptData[index][0];
            sendScriptCheckVersion(name);
        }

        void AppConfig_Click(object sender, System.EventArgs e)
        {
            string path = getPluginDirectory() + "\\config\\_app.yaml";
            Process.Start("notepad.exe", path);
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
            } catch (System.Net.Sockets.SocketException) {
            }
            // listen from daemon
            if (udpListener != null) {
                udpListener.Close();
            }
            udpListener = new UdpClient();
            listenThread = new Thread(new ThreadStart(udpListen));
            listenThread.Start();
            sendScriptRequest();
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
                        case DATA_TYPE_SCRIPT: {
                            string valStr = System.Text.Encoding.UTF8.GetString(bytes, 1, bytes.Length-1);                           
                            var valSplit = valStr.Split('|');
                            var hasScript = false;
                            for (var i = 0; i < this.scriptData.Count; i++) {
                                var item = this.scriptData[i];
                                if (item[0] == valSplit[0]) {
                                    hasScript = true;
                                    this.scriptData[i] = valSplit;
                                    this.formScriptList.Items[i] = "[" + (valSplit[1] == "" ? " " : "O") + "] " + valSplit[2];
                                    break;
                                }
                            }
                            foreach (var item in this.scriptData) {
                                if (item[0] == valSplit[0]) {
                                    hasScript = true;
                                    break;
                                }
                            }
                            if (hasScript) {
                                break;
                            }
                            this.scriptData.Add(valSplit);
                            string label = "[" + (valSplit[1] == "" ? " " : "O") + "] " + valSplit[2];
                            this.formScriptList.Items.Add(label);
                            if (this.lastScriptSelected != "" && valSplit[0] == this.lastScriptSelected) {
                                int index = this.formScriptList.FindString(label);
                                this.formScriptList.SetSelected(index, true);
                            }
                            break;
                        }
                        case DATA_TYPE_ACT_UPDATE: {
                            DialogResult result = MessageBox.Show("An updated version of the FFTools plugin is ready. Please restart ACT for the changes to take effect.", "New Version", MessageBoxButtons.OK, MessageBoxIcon.Asterisk);
                            break;
                        }
                        case DATA_TYPE_SCRIPT_VERSION: {
                            string valStr = System.Text.Encoding.UTF8.GetString(bytes, 1, bytes.Length-1);                           
                            var valSplit = valStr.Split('|');
                            var name = valSplit[0];
                            var currentVersion = valSplit[1];
                            var latestVersion = valSplit[2];
                            if (currentVersion != latestVersion) {
                                DialogResult result = MessageBox.Show("An update is available. Update now? (v" + currentVersion + " => v" + latestVersion + ")", "Script Version Check", MessageBoxButtons.YesNo, MessageBoxIcon.Information);
                                if (result == DialogResult.Yes) {
                                    sendScriptUpdateVersion(name);
                                }
                                break;
                            }
                            MessageBox.Show("You have the latest version. (v" + latestVersion + ")", "Script Version Check", MessageBoxButtons.OK, MessageBoxIcon.Information);
                            break;
                        }
                        case DATA_TYPE_SCRIPT_UPDATE: {
                            string valStr = System.Text.Encoding.UTF8.GetString(bytes, 1, bytes.Length-1);                           
                            var valSplit = valStr.Split('|');
                            var name = valSplit[0];
                            var status = valSplit[1];
                            switch (status) {
                                case "success": {
                                    MessageBox.Show("Update complete.", "Script Update", MessageBoxButtons.OK, MessageBoxIcon.Information);
                                    break;
                                }
                                default: {
                                    MessageBox.Show("An error occured during script update.\n\n" + valSplit[2], "Script Update", MessageBoxButtons.OK, MessageBoxIcon.Exclamation);
                                    break;
                                }
                            }
                            break;
                        }
                    }

                } catch (SocketException) {
                }
            }
            //udpListener.Close();
        }

        void listenDaemonExit()
        {
            if (this.scriptDaemon == null) {
                return;
            }
            string err = this.scriptDaemon.StandardError.ReadToEnd();
            this.scriptDaemon.WaitForExit();
            MessageBox.Show("The FFTools Lua daemon closed unexpectedly (exit code " + this.scriptDaemon.ExitCode + ")\n\n"+err, "FFTools Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
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
            if (actionInfo.cancelAction || !actionInfo.tags.ContainsKey("SourceId")) {
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
                        (string) actionInfo.tags["SourceId"],
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
            prepareInt32(ref sendData, cd.CritHits);                       // number of critical hits
            prepareInt32(ref sendData, cd.CritHeals);                      // number of critical heals
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

        void resetScriptList()
        {
            if (this.scriptData != null && this.scriptData.Count > 0 && this.formScriptList.SelectedItem != null) {
                string value = this.formScriptList.SelectedItem.ToString();
                int index = this.formScriptList.FindString(value);
                this.lastScriptSelected = this.scriptData[index][0];
            }
            this.scriptData = new List<string[]>();
            this.formScriptList.Items.Clear();
            this.formScriptList.Enabled = true;
        }

        void sendScriptRequest()
        {
            this.resetScriptList();
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_SCRIPT);
            sendUdp(ref sendData);
        }

        void sendEnableScript(string name)
        {
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_SCRIPT_ENABLE);
            prepareString(ref sendData, name);
            sendUdp(ref sendData);
        }

        void sendDisableScript(string name)
        {
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_SCRIPT_DISABLE);
            prepareString(ref sendData, name);
            sendUdp(ref sendData);
        }

        void sendReloadScript(string name)
        {
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_SCRIPT_RELOAD);
            prepareString(ref sendData, name);
            sendUdp(ref sendData);
        }

        void sendScriptCheckVersion(string name)
        {
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_SCRIPT_VERSION);
            prepareString(ref sendData, name);
            sendUdp(ref sendData);
        }

       void sendScriptUpdateVersion(string name)
        {
            List<Byte> sendData = new List<Byte>();
            sendData.Add(DATA_TYPE_SCRIPT_UPDATE);
            prepareString(ref sendData, name);
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

        string getPluginDirectory()
        {
            foreach (ActPluginData p in ActGlobals.oFormActMain.ActPlugins) {
                if (p.pluginObj == this) {
                    return p.pluginFile.DirectoryName;
                }
            }
            return "";
        }

        string getScriptConfigFile(string name)
        {
            string path = getPluginDirectory() + "\\config\\" + name + ".yaml";
            if (!File.Exists(path)) {
                return "";
            }
            return path;
        }

        void onScriptConfigChanges(object source, FileSystemEventArgs e)
        {
            string scriptName = Path.GetFileNameWithoutExtension(e.Name);
            sendReloadScript(scriptName);
        }

    }
}
