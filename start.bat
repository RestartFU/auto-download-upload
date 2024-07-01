@echo off
set servers=play.craftia.me:19132 realmc.xyz:19132 tetramc.org:19132 soulmine.ru:19132 skill-mine.ru:19133 play.rustmc.online:19132 geo.hivebedrock.network:19132
set len=0

(for %%s in (%servers%) do (
    set /a len+=1
))

 :loop
    (for %%s in (%servers%) do (
       start /B /wait downloader.exe %%s decryptmypack.com
        timeout /t 600 /nobreak >nul
        echo[
    ))
 goto loop