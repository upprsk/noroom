let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/dev/upprsk/noroom/app
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +16 src/routes/+page.server.ts
badd +1 ~/dev/upprsk/noroom/app/src/routes/+layout.server.ts
badd +27 ~/dev/upprsk/noroom/app/src/routes/+layout.svelte
badd +110 src/lib/models.ts
badd +16 src/lib/components/BasicFormCard.svelte
badd +70 src/lib/pocketbase.ts
badd +22 src/lib/components/NavbarEnd.svelte
badd +21 ~/dev/upprsk/noroom/app/src/routes/+page.svelte
badd +38 ~/dev/upprsk/noroom/app/src/routes/classes/new/+page.svelte
badd +23 src/routes/register/+page.svelte
badd +1 ~/dev/upprsk/noroom/app/src/routes/classes/new/+page.server.ts
badd +19 ~/dev/upprsk/noroom/app/src/routes/register/+page.server.ts
badd +3 src/lib/components/BasicAvatar.svelte
badd +1 src/lib/components/BasicCard.svelte
badd +1 ~/dev/upprsk/noroom/app/src/lib/components/input/TextInput.svelte
badd +43 ~/dev/upprsk/noroom/app/src/lib/components/input/TextArea.svelte
badd +80 src/routes/classes/\[id]/+page.svelte
badd +18 ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/+page.server.ts
badd +16 ~/dev/upprsk/noroom/app/eslint.config.js
badd +84 ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte
badd +65 ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.server.ts
badd +3 ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/FileUploadDialog.svelte
argglobal
%argdel
set stal=2
tabnew +setlocal\ bufhidden=wipe
tabrewind
edit ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/FileUploadDialog.svelte
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd _ | wincmd |
split
1wincmd k
wincmd w
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe '1resize ' . ((&lines * 32 + 34) / 68)
exe 'vert 1resize ' . ((&columns * 159 + 159) / 319)
exe '2resize ' . ((&lines * 32 + 34) / 68)
exe 'vert 2resize ' . ((&columns * 159 + 159) / 319)
exe 'vert 3resize ' . ((&columns * 159 + 159) / 319)
argglobal
balt ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 3 - ((2 * winheight(0) + 16) / 32)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 3
normal! 0
wincmd w
argglobal
if bufexists(fnamemodify("~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.server.ts", ":p")) | buffer ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.server.ts | else | edit ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.server.ts | endif
if &buftype ==# 'terminal'
  silent file ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.server.ts
endif
balt ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/FileUploadDialog.svelte
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 63 - ((27 * winheight(0) + 16) / 32)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 63
normal! 07|
wincmd w
argglobal
if bufexists(fnamemodify("~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte", ":p")) | buffer ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte | else | edit ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte | endif
if &buftype ==# 'terminal'
  silent file ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte
endif
balt ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.server.ts
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 83 - ((61 * winheight(0) + 32) / 65)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 83
normal! 0
wincmd w
2wincmd w
exe '1resize ' . ((&lines * 32 + 34) / 68)
exe 'vert 1resize ' . ((&columns * 159 + 159) / 319)
exe '2resize ' . ((&lines * 32 + 34) / 68)
exe 'vert 2resize ' . ((&columns * 159 + 159) / 319)
exe 'vert 3resize ' . ((&columns * 159 + 159) / 319)
tabnext
edit src/lib/models.ts
argglobal
balt ~/dev/upprsk/noroom/app/src/routes/classes/\[id]/edit/+page.svelte
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 107 - ((35 * winheight(0) + 32) / 65)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 107
normal! 010|
tabnext 1
set stal=1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
