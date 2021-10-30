local encounter = {}
local combatants = {}

local function on_encounter(e)
    encounter = e
end

local function on_encounter_change(e)
    combatants = {}
end

local function on_combatant(c)
    combatants[c.id] = c
end

function init()
    event_attach("act:encounter", on_encounter)
    event_attach("act:encounter:change", on_encounter_change)
    event_attach("act:combatant", on_combatant)
end

function web(req)
    return {
        encounter = encounter,
        combatants = combatants,
        me = me()
    }
end

function info()
    return {
        name = "Live Parse",
        desc = "Output parses to web server in real time."
    }
end