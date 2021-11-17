local jail_list = {}
local encounter_id = 0
local has_called = false

local function jail_sort(a, b)
    local order_config = fft_config_get("order")
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
        local ca = fft_combatant_from_name(a)
        local cb = fft_combatant_from_name(b)
        if ak == -1 and ca ~= nil and string.upper(ca.job) == string.upper(v) then
            ak = k
        end
        if bk == -1 and cb ~= nil and string.upper(cb.job) == string.upper(v) then
            bk = k
        end
    end
    return ak < bk
end

local function clear_marks()
    fft_event_dispatch("am:clear")
end

local function on_log(l)
    local match = fft_regex_match(":2B6(B|C):.*?:.*?:(.*?):0:", l.log_line)
    if match == nil then
        return
    end
    if #jail_list < 3 then
        jail_list[#jail_list+1] = match[3]
    end
    if #jail_list == 3 and not has_called then
        table.sort(jail_list, jail_sort)
        for n, v in ipairs(jail_list) do
            fft_say_if(n, {name=v})
            fft_event_dispatch("am:mark", v)
        end
        fft_wait(8000, clear_marks)
        has_called = true
    end
end

local function on_encounter(e)
    jail_list = {}
    encounter_id = e.id
    has_called = false
end

function web()
    return {
        encounter_id = encounter_id,
        jails = jail_list
    }
end

function init()
    fft_event_attach("act:log_line", on_log)
    fft_event_attach("act:encounter:change", on_encounter)
end

function info()
    return {
        name = "UWU Jails",
        desc = "Calls out your Titan jail number in The Weapon's Refrain (Ultimate). If the 'Auto Markers' script is enabled then this will also mark the players.",
        version = "0.01"
    }
end