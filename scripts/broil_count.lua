broilCount = 0

local function on_log(l)
    if l.ability_id == 0x409D and me().id == l.source_id then
        broilCount = broilCount + 1
        data_set("broil_count", broilCount)
        act_say(broilCount)
    end
end

function init()
    event_attach("act:log_line", on_log)
    broilCount = data_get("broil_count")
    if broilCount == nil then
        broilCount = 0
    end
end

function info()
    return {
        name = "Broil Counter",
        desc = "Counters the number of Broil casts."
    }
end