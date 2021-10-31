function init()
end

function web(req)
    return {
        encounter = act_encounter(),
        combatants = act_combatants(),
        me = me()
    }
end

function info()
    return {
        name = "Live Parse",
        desc = "Output parses to web server in real time."
    }
end