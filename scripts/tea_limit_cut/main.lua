local function on_log(l)
    if l.type == 0x1B and l.icon_id >= 79 and l.icon_id <= 86 then
        local number = l.icon_id - 78
        fft_say_if(number, {name=l.target_name})
    end
end

function init()
    fft_event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "TEA Limit Cut Number",
        desc = "Calls the player's limit cut number in The Epic of Alexander (Ultimate) and other fights that use the same markers.",
        version = "0.01"
    }
end