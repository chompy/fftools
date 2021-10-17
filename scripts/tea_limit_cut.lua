local function on_log(l)
    local matches = regex_match("1B:.*:(.*):0000:.*:(00[4-5][F0-6]):", l.log_line)
    if matches == nil or matches[2] == nil then
        return
    end
    if matches[2] == me().name then
        local number = tonumber(matches[3], 16) - 78
        if number > 1 and number < 9 then
            act_say(number)
        end
    end

end

function init()
    event_attach("act:log_line", on_log)
end