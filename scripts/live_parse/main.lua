function init()
end

function web(req)
    return {
        encounter = fft_encounter(),
        combatants = fft_combatants(),
        me = fft_me()
    }
end

function info()
    return {
        name = "Live Parse",
        desc = "Output parses to web server in real time."
    }
end
