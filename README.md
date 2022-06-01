FFTools (Final Fantasy XIV Tools)
================================
**By Chompy / Minda Silva@Sargatanas / Qunara Sivra@Excalibur**

Extends Final Fantasy XIV log parsing in Advanced Combat Tracker (ACT) with Lua scripts that can perform TTS callouts, create web UIs, and more.


## Installation

1. Download the plugin [here](https://github.com/chompy/fftools/releases/latest). Extract the ZIP. Open ACT and navigate to the plugins tab.
2. Click 'Browse...' and locate the 'FFTools_ACT_Plugin.cs' file.
3. Click 'Add/Enable Plugin.'
4. Click on the 'FFTools' tab. A list of available scripts will be on the left side. Click on a script name and then click the 'Enable' button to enable the script.

- You might get a Windows firewall alert. This is because the ACT plugin launches a seperate application to run the Lua scripts. This application communicates with ACT over an internal network connection.

## Web View

Some scripts provide a web view which provides additonal visual information. These web views can be used in OBS as part of your streaming overlay or can be shared with other players who can't use ACT. To access the web view just click "Open Web View" while the desired script is selected, it should then open a new browser tab. The link in the browser's address bar can be shared with other players!

### Disable Sharable Web Views

By default web views are made publically available via a proxy to fftools.net. This can be disabled by clicking "Edit Main Plugin Config" and changing the line (in the resulting notepad file that opens) containing `enable_proxy: true` to `enable_proxy: false`.
Web views will still be available on your local machine at http://localhost:31594. You can also enable port forwarding in your router to share web views without the use of the FFTools proxy.


## Scripting API

Scripts are expected to contain two global functions, `init` and `info`. The `info` function is called to obtain information about the script, a name and a description. It should return a table with keys `name` and `desc`. The `init` function is called when the script is first enabled. It is expected to attach to any event needed by the script.

Example script...

```
local function on_zone(z)
    print("Enter zone " .. z)
end

function info()
    return {
        name = "Example Script",
        desc = "A very basic example script."
    }
end

function init()
    fft_event_attach("act:encounter:zone", on_zone)
end
```


### Available Functions

Below is a list of available functions that can be used in your scripts. Better documentation will be provided later.

- fft_event_attach
- fft_event_detach
- fft_event_dispatch
- fft_say
- fft_say_if
- fft_combatants
- fft_combatant_from_id
- fft_combatant_from_name
- fft_config_get
- fft_data_set
- fft_data_get
- fft_log_info
- fft_log_warn
- fft_me
- fft_regex_match
- fft_wait
- fft_key_press


### Available Events

- act:log_line
- act:combatant
- act:encounter
- act:encounter:zone
- act:encounter:change

