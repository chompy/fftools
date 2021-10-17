broilCount = 0

local function log_line(l)
    matches = regex_match("15:.*:(.*):409D:Broil III:", l.log_line)
    if matches == nil or matches[2] == nil then
        return
    end
    if me().name == matches[2] then
        broilCount = broilCount + 1
        act_say(broilCount)
    end
end

function init()
    event_attach("act:log_line", log_line)
end