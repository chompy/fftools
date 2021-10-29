local function on_log(l)
    if l.type == 0x1B and l.icon_id >= 79 and l.icon_id <= 86 and l.target_id == me().id then
        local number = l.icon_id - 78
        act_say(number)
    end
end

function init()
    event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "Limit Cut Number",
        desc = "Calls the player's limit cut number."
    }
end