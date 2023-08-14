## Rtorrent Flood

login to website create user account. then exec into the rtorrent pod and create two folders

set the download folder to /media/download

lets now create a few folders on the shared-media pv for TV/MOVIES
```
❯ kubectl exec -it -n media rtorrent-flood-785b55896c-hjp86 ash
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
/ $ cd /media
/media $ mkdir TV
/media $ mkdir MOVIES

```
path in setup: ~/.local/share/rtorrent/rtorrent.sock

set download path in settings to: /media/download