package adif

import (
  "io"
  "log"
)


%%{
  machine formula;
  write data;

  langle = '<';
  rangle = '>';
  colon = ':';
  nl = '\n' | '\r';
  eoh = /<eoh>/i;
  eor = /<eor>/i;
  
  
  main := |*
  nl { l.emit(tokenNL,data[ts:te]) };
  colon { l.emit(tokenColon, data[ts:te])};
  langle { l.emit(tokenLAngle, data[ts:te])};
  rangle { l.emit(tokenRAngle, data[ts:te])};
  eoh { l.emit(tokenEOH,data[ts:te]) };
  eor { l.emit(tokenEOR,data[ts:te]) };
  any { l.emit(tokenOther, data[ts:te])};
*|;


prepush {
  stack = append(stack,0)
}

postpop {
  stack = stack[0:len(stack)-1]
}

}%%
 func(l *Lexer) lex(r io.Reader)  {
  cs, p, pe := 0, 0, 0
  eof := -1
  ts, te,act := 0,0,0
  _ = act
  curline := 1
  _ = curline
  data := make([]byte,4096)

done := false
for !done {
 // p - index of next character to process
 // pe - index of the end of the data
 // eof - index of the end of the file
 // ts - index of the start of the current token
 // te - index of the end of the current token

   // still have a partial token 
   rem := 0
   if ts > 0 {
      rem = p - ts
   }
   p = 0
   n, err := r.Read(data[rem:])
   if err == io.EOF {
     l.emit(tokenEOF,nil)
   }
   if n == 0 || err != nil {
     done = true
   }
   pe = n+rem
   if pe < len(data) {
      eof = pe
   }

  %%{
    write init;
    write exec;
  }%%
  
  if ts > 0 {
      // currently parsing a token, so shift it to the
      // beginning of the buffer
      copy(data[0:],data[ts:])
    } 
  }
  
  if cs == formula_error {
     l.emit(tokenLexError,nil)
  }
  close(l.nodes)
}
