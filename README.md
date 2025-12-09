# Local Mail
It is a reliable messenger for local network that implemets CIA model. Created using golang and fyne framework
<br><br>
## Features
Messenger provides p2p mTLS text/file messages to specific user in network, that is reached via udp broadcast.
<table width="100%">
  <tr>
    <td width="50%" align="center">
      <img src="https://github.com/user-attachments/assets/43a05a80-88e0-4ce6-b931-6dd556e99136" width="100%">
    </td>
    <td width="50%" align="center">
      <img src="https://github.com/user-attachments/assets/b169b640-c58d-4509-982e-208bb30bf70f" width="100%">
    </td>
  </tr>
</table>
<br>

## Screenshots
<table width="100%">
  <tr>
    <td width="50%" align="center">
      <img src="https://github.com/user-attachments/assets/d5078eb9-b5a1-4502-9c09-05b58adf02c2" width="100%">
      <br><b>Refreshed chat list</b>
    </td>
    <td width="50%" align="center">
      <img src="https://github.com/user-attachments/assets/5ee9d69f-1e83-425c-a230-2804ead8ec5f" width="100%">
      <br><b>Chat example</b>
    </td>
  </tr>
  <tr>
    <td width="50%" align="center">
      <img src="https://github.com/user-attachments/assets/25895f45-dd5d-4db8-8a68-61613b8a6781" width="100%">
      <br><b>Menu</b>
    </td>
    <td width="50%" align="center">
      <img src="https://github.com/user-attachments/assets/617cb0ce-21f6-4009-b9d7-aebc78e37959" width="100%">
      <br><b>Settings</b>
    </td>
  </tr>
</table>
<br>

## Main questions
<hr>
<h3>Messanger is not working after changing ports/installation.</h3><br>
After changing port you have to allow your firewall this port/protocol.<br>

#### .deb
```bash
sudo ufw allow 1337/udp && sudo ufw allow 1338/tcp
```
#### .rpm 
```bash
sudo firewall-cmd --permanent --add-port=1338/tcp
sudo firewall-cmd --permanent --add-port=1337/udp
sudo firewall-cmd --reload
```

#### .exe
```powershell
New-NetFirewallRule -DisplayName "App UDP 1337" -Direction Inbound -LocalPort 1337 -Protocol UDP -Action Allow
New-NetFirewallRule -DisplayName "App TLS 1338" -Direction Inbound -LocalPort 1338 -Protocol TCP -Action Allow
```

#### .apk (ADB Shell or Terminal)
```bash
su -c 'iptables -A INPUT -p udp --dport 1337 -j ACCEPT'
su -c 'iptables -A INPUT -p tcp --dport 1338 -j ACCEPT'
```
<hr>
