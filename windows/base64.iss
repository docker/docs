// Source: http://www.vincenzo.net/isxkb/index.php?title=Encode/Decode_Base64
const Codes64 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';

function Encode64(S: AnsiString): AnsiString;
var
	i: Integer;
	a: Integer;
	x: Integer;
	b: Integer;
begin
	Result := '';
	a := 0;
	b := 0;
	for i := 1 to Length(s) do
	begin
		x := Ord(s[i]);
		b := b * 256 + x;
		a := a + 8;
		while (a >= 6) do
		begin
			a := a - 6;
			x := b div (1 shl a);
			b := b mod (1 shl a);
			Result := Result + copy(Codes64,x + 1,1);
		end;
	end;
	if a > 0 then
	begin
		x := b shl (6 - a);
		Result := Result + copy(Codes64,x + 1,1);
	end;
	a := Length(Result) mod 4;
	if a = 2 then
		Result := Result + '=='
	else if a = 3 then
		Result := Result + '=';

end;

function Decode64(S: AnsiString): AnsiString;
var
	i: Integer;
	a: Integer;
	x: Integer;
	b: Integer;
begin
	Result := '';
	a := 0;
	b := 0;
	for i := 1 to Length(s) do
	begin
		x := Pos(s[i], codes64) - 1;
		if x >= 0 then
		begin
			b := b * 64 + x;
			a := a + 6;
			if a >= 8 then
			begin
				a := a - 8;
				x := b shr a;
				b := b mod (1 shl a);
				x := x mod 256;
				Result := Result + chr(x);
			end;
		end
	else
		Exit; // finish at unknown
	end;
end;
