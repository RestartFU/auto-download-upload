@echo off
set servers=geo.hivebedrock.network:19132 play.rustmc.online:19132

(for %%s in (%servers%) do (
   start /wait downloader.exe %%s decryptmypack.com
))