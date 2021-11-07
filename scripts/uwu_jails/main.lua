local jail_list = {}
local encounter_id = 0

local function jail_sort(a, b)
    local order_config = ffl_config_get("order")
    local ak = -1
    local bk = -1
    -- match name
    for k, v in ipairs(order_config) do
        if ak == -1 and a == v then
            ak = k
        end
        if bk == -1 and b == v then
            bk = k
        end
    end
    -- match job
    for k, v in ipairs(order_config) do
        local ca = ffl_combatant_from_name(a)
        local cb = ffl_combatant_from_name(b)
        if ak == -1 and ca ~= nil and string.upper(ca.job) == string.upper(v) then
            ak = k
        end
        if bk == -1 and cb ~= nil and string.upper(cb.job) == string.upper(v) then
            bk = k
        end
    end
    return ak < bk
end

local function on_log(l)
    local match = ffl_regex_match(":2B6(B|C):.*?:.*?:(.*?):0:", l.log_line)
    if match == nil then
        return
    end
    if #jail_list < 3 then
        jail_list[#jail_list+1] = match[3]
    end
    if #jail_list == 3 then
        table.sort(jail_list, jail_sort)
        for n, v in ipairs(jail_list) do
            ffl_say_if(n, {name=v})
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
    ffl_event_attach("act:log_line", on_log)
    ffl_event_attach("act:encounter:change", on_encounter)
end

function info()
    return {
        name = "UWU Jails Callout",
        desc = "Calls out your Titan jail number."
    }
end