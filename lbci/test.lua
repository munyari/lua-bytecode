local inspector=require"bci"

local function inspect(f,all)
 local F=inspector.getheader(f)
 print("header",f)
 for k,v in next,F do
  print("",k,v)
 end

 print("constants",F.constants)
 for i=1,F.constants do
  local a=inspector.getconstant(f,i)
  print("",i,a)
 end

 print("locals",F.locals)
 for i=1,F.locals do
  local a,b,c=inspector.getlocal(f,i)
  print("",i,a,b,c,i<=F.params and "param" or "")
 end

 print("upvalues",F.upvalues)
 for i=1,F.upvalues do
  local a,b,c=inspector.getupvalue(f,i)
  print("",i,a,b,c)
 end

 print("functions",F.functions)
 for i=1,F.functions do
  local a=inspector.getfunction(f,i)
  print("",i,a)
 end

 print("instructions",F.instructions)
 for i=1,F.instructions do
  local a,b,c,d,e=inspector.getinstruction(f,i)
  local b=string.sub(b.."          ",1,9)
  print("",i,a,b,c,d,e)
 end

 print()
 if all then
  for i=1,F.functions do
   inspect(inspector.getfunction(f,i),all)
  end
 end

end

local function globals(f,all)
 local F=inspector.getheader(f)
 for i=1,F.instructions do
  local a,b,c,d,e=inspector.getinstruction(f,i)
  if b=="GETTABUP" and inspector.getupvalue(f,d+1)=="_ENV" then
   print("",F.source,a,"GET ",inspector.getconstant(f,-e))
  elseif b=="SETTABUP" and inspector.getupvalue(f,c+1)=="_ENV" then
   print("",F.source,a,"SET*",inspector.getconstant(f,-d))
  elseif b=="GETTABLE" and inspector.getlocal(f,d+1)=="_ENV" then
   print("",F.source,a,"GET ",inspector.getconstant(f,-e))
  elseif b=="SETTABLE" and inspector.getlocal(f,c+1)=="_ENV" then
   print("",F.source,a,"SET*",inspector.getconstant(f,-d))
  end
 end
 if all then
  for i=1,F.functions do
   globals(inspector.getfunction(f,i),all)
  end
 end
end

local f=assert(loadfile(arg[1] or "sample.lua"))
inspect(f,true)

print"globals"
globals(f,true)

print"setconstant"
f=function() print("","hello") end	f()
inspector.setconstant(f,3,"bye")	f()
inspector.setconstant(f,3)		f()
inspector.setconstant(f,3,1993)		f()
inspector.setconstant(f,3,print)	f()
inspector.setconstant(f,3,math)		f()
