// Licensed under Creative Commons Attribution-ShareAlike 3.0 Unported
// Source: http://stackoverflow.com/questions/1347083/how-to-generate-a-new-guid-in-inno-setup

function CoCreateGuid(var Guid:TGuid):integer;
external 'CoCreateGuid@ole32.dll stdcall';

function inttohex(l:longword; digits:integer):string;
var hexchars:string;
begin
 hexchars:='0123456789ABCDEF';
 setlength(result,digits);
 while (digits>0) do begin
  result[digits]:=hexchars[l mod 16+1];
  l:=l div 16;
  digits:=digits-1;
 end;
end;

function GetGuid(dummy:string):string;
var Guid:TGuid;
begin
  if CoCreateGuid(Guid)=0 then begin
  result:='{'+IntToHex(Guid.D1,8)+'-'+
           IntToHex(Guid.D2,4)+'-'+
           IntToHex(Guid.D3,4)+'-'+
           IntToHex(Guid.D4[0],2)+IntToHex(Guid.D4[1],2)+'-'+
           IntToHex(Guid.D4[2],2)+IntToHex(Guid.D4[3],2)+
           IntToHex(Guid.D4[4],2)+IntToHex(Guid.D4[5],2)+
           IntToHex(Guid.D4[6],2)+IntToHex(Guid.D4[7],2)+
           '}';
  end else
    result:='{00000000-0000-0000-0000-000000000000}';
end;
