-- function foo()
--     return 3 + 5
-- end
-- function foo()
--     return 234342
-- end

-- a = foo()
-- b = 1203223
-- print(a + b)

print(3 + 5)


-- 10101011011 | 00

-- 00 + (higher_byte & 1) >> 2
--[[ OP_ADD, /*	A B C	R(A) := RK(B) + RK(C)				*/
 OP_ADD is 13
 13 = 8 + 4 + 1 = 0b001101
 d in hex
6 bits for d
8 bits for A
9 bits for B
9 bits for C

]]--

