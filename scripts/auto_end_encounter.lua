local function on_log(l)
    local match = regex_match("21:.*:(40000001|40000005|40000003):00:00:00:00", l.log_line)
    if match ~= nil and match[1] ~= nil then
        act_end()
    end
end

function init()
    event_attach("act:log_line", on_log)
end