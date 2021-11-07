function init()
end

function web(req)
    return {
        encounter = ffl_encounter(),
        combatants = ffl_combatants(),
        me = ffl_me()
    }
end

function info()
    return {
        name = "Live Parse",
        desc = "Output parses to web server in real time."
    }
end
