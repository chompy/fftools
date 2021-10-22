local is_tea_zone = false

local function on_zone(z)
    is_tea_zone = false
    if z == "The Epic Of Alexander (Ultimate)" then
        is_tea_zone = true
        log_info("Entered TEA.")
    end
end

local function on_log(l)
    if not is_tea_zone then
        return
    end
    local matches = regex_match("1B:.*:(.*):0000:.*:(00[4-5][F0-6]):", l.log_line)
    if matches == nil or matches[2] == nil then
        return
    end
    if matches[2] == me().name then
        local number = tonumber(matches[3], 16) - 78
        if number > 0 and number < 9 then
            act_say(number)
        end
    end
end

function init()
    event_attach("act:encounter:zone", on_zone)
    event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "[TEA] Limit Cut Number",
        desc = "Calls the player's limit cut number in The Epic Of Alexander (Ultimate)."
    }
end