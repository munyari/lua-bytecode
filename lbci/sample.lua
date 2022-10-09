-- nonsense program to exercise VM instructions

-- test globals
do
  --do local x,y,z end
  local _ENV = {}
  _NAME = 'BAZ'
  function Foo() if _NAME then return end end
  function Bar() Foo() end
end

local x
function f(a,...) local b=... x=89-x return a*x+b,... end
local t={100,200,300}
for i=1,3 do t[x..i]=-i^2/10 print(t[i],#t) end
for k,v in pairs(t) do print(k,v) end
return x<9, not x, x:run("down")
