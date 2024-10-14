local function put(...)
	io.write(...)
	io.write "\n"
end

local function cleanup(rems)
	os.remove "tmp.exe";
	if rems then os.remove "tmp.s" end
	os.remove "9cc.exe" -- ???
end

cleanup(true)

local function cmd(c, exit)
	local res, term, code = os.execute(c)
	if not res and (exit or term ~= "exit") then
		put("\ncommand '", c, "' terminated abnormally with signal/code ", code)
		cleanup()
		os.exit(false)
	end
	return code
end

put "compiling executable..."
cmd("go build -o 9cc.exe .", true)
put "done\n" -- extra \n intended

local function test(expected, input)
	cmd('.\\9cc.exe "'..input..'" > tmp.s', true)
	cmd "cc -o tmp tmp.s"
	local actual = cmd ".\\tmp.exe"
	local res = actual == expected

	io.write("'", tostring(input), "' => ")
	if not res then
		io.write(tostring(expected), " expected, but got ")
	end
	put(tostring(actual))

	if not res then
		cleanup()
		os.exit(1)
	end
end

put "testing..."

-- single numbers
test(0,  "0" )
test(42, "42")

-- basic tokenization
test(21, "5+20-4"       )
test(41, " 12 + 34 - 5 ")

-- complex expressions
test(47, "5+6*7"        )
test(15, "5*(9-6)"      )
test(4,  "(3+5)/2"      )

-- unary expressions
test(-5,  "-5"    )
test(20,  "--20"  )
test(20,  "--+20" )
test(-10, "10-20" )
test(10,  "-10+20")

put "OK"
cleanup(true)
