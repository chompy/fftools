local eventLogLine = nil

local function test(c)
    print(c.damage)
    --event_detach(eventLogLine)
end

function init()
    eventLogLine = event_attach("act:combatant", test)
end

