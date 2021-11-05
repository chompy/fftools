local jail_list = {}
local encounter_id = 0

local function jail_sort(a, b)
    local order_config = config_get("order")
    local ak = -1
    local bk = -1
    for k, v in ipairs(order_config) do
        local ca = act_combatant_from_name(a)
        local cb = act_combatant_from_name(b)
        if ak == -1 and (a == v or (ca ~= nil and string.upper(ca.job) == string.upper(v))) then
            ak = k
        end
        if bk == -1 and (b.name == v or (cb ~= nil and string.upper(cb.job) == string.upper(v))) then
            bk = k
        end
    end
    return ak < bk
end

local function on_log(l)
    local match = regex_match(":2B6(B|C):.*?:.*?:(.*?):0:", l.log_line)
    if match == nil then
        return
    end
    if #jail_list < 3 then
        jail_list[#jail_list+1] = match[3]
    end
    if #jail_list == 3 then
        table.sort(jail_list, jail_sort)
        for n, v in ipairs(jail_list) do
            if v == me().name then
                act_say(n)
                return
            end
        end
    end
end

local function on_encounter(e)
    jail_list = {}
    encounter_id = e.id
end

function web()
    return {
        encounter_id = encounter_id,
        jails = jail_list
    }
end

function init()
    event_attach("act:log_line", on_log)
    event_attach("act:encounter:change", on_encounter)
end

function info()
    return {
        name = "UWU Jails Callout",
        desc = "Calls out your Titan jail number."
    }
end