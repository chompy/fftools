FFTK (Final Fantasy XIV Toolkit)
================================
**By Chompy / Minda Silva@Sargatanas / Qunara Sivra@Excalibur**

Extends Final Fantasy XIV log parsing in Advanced Combat Tracker (ACT) with Lua scripts that can perform TTS callouts, create web UIs, and more.


## Installation

1. Download the plugin here. Extract the ZIP. Open ACT and navigate to the plugins tab.
2. Click 'Browse...' and locate the 'FFXIV_ACT_Lua.cs' file.
3. Click 'Add/Enable Plugin.'
4. Click on the 'FFXIV Lua Scripts' tab. A list of available scripts will be on the left side. Click on a script name and then click the 'Enable' button to enable the script.


## Web View

Some scripts provide a webpage view which provides additonal visual information. These web views can be used in OBS as part of your streaming overlay or can be shared with other players who can't use ACT.

When using 


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
    ffl_event_attach("act:encounter:zone", on_zone)
end
```


### Available Functions

Below is a list of available functions that can be used in your scripts. Better documentation will be provided later.

- ffl_event_attach
- ffl_event_detach
- ffl_event_dispatch
- ffl_say
- ffl_say_if
- ffl_combatants
- ffl_combatant_from_id
- ffl_combatant_from_name
- ffl_config_get
- ffl_data_set
- ffl_data_get
- ffl_log_info
- ffl_log_warn
- ffl_me
- ffl_regex_match
- ffl_wait
- ffl_key_press


### Available Events

- act:log_line
- act:combatant
- act:encounter
- act:encounter:zone
- act:encounter:change

