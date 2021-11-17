local function on_log(l)
    if string.match(l.log_line, "Blazing path, lead me to iron rule!") then
        fft_say("stack and out")
    elseif string.match(l.log_line, "Fleeting light! Amid a rain of stars, exalt you the red moon!") then
        fft_say("spread and dive")
    elseif string.match(l.log_line, "Fleeting light! 'Neath the red moon, scorch you the earth!") then
        fft_say("dive and stack")
    elseif string.match(l.log_line, "From on high I descend, the hallowed moon to call!") then
        fft_say("dive and in")
    elseif string.match(l.log_line, "From on high I descend, the iron path to call!") then
        fft_say("dive and out")
    elseif string.match(l.log_line, "From on high I descend, the iron path to walk!") then
        fft_say("dive and out")
    elseif string.match(l.log_line, "O hallowed moon, shine you the iron path!") then
        fft_say("in and out")
    elseif string.match(l.log_line, "O hallowed moon, take fire and scorch my foes!") then
        fft_say("in and stack")
    elseif string.match(l.log_line, "Take fire, O hallowed moon!") then
        fft_say("stack and in")
    elseif string.match(l.log_line, "From hallowed moon I descend, a rain of stars to bring!") then
        fft_say("in, dive, and spread")
    elseif string.match(l.log_line, "From on high I descend, the moon and stars to bring!") then
        fft_say("dive, in, spread")
    elseif string.match(l.log_line, "From hallowed moon I bare iron, in my descent to wield!") then
        fft_say("in, out, dive")
    elseif string.match(l.log_line, "Unbending iron, descend with fiery edge!") then
        fft_say("out, dive, stack")
    elseif string.match(l.log_line, "From hallowed moon I descend, upon burning earth to tread!") then
        fft_say("in, dive, stack")
    elseif string.match(l.log_line, "Unbending iron, take fire and descend!") then
        fft_say("out, stack, dive")
    end
end

function init()
    fft_event_attach("act:log_line", on_log)
end

function info()
    return {
        name = "UCOB Nael Quotes",
        desc = "Calls out the Nael quote mechanics in The Unending Coil of Bahamut (Ultimate).",
        version = "0.01"
    }
end