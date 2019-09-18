# fritzbox-cloudflare-dnsupdater

`fritzbox-cloudflare-dnsupdater` is meant to be used as a standalone webservice that allows to remotely edit DNS
 Address records at [Cloudflare](https://www.cloudflare.com/).  
 
 It was developed to address the special needs when populating IP address changes received at a [AVM FRITZ!Box](https://avm.de/produkte/fritzbox/) to Cloudflare.
 
## Docker Image

A Docker image is available from official Docker registry at `suckowbiz/fritzbox-cloudflare-dnsupdater`.

## Usage

1. Start `fritzbox-cloudflare-dnsupdater` (e.g. on a machine in the local AVM FRITZ!Box network:

    ```bash
    $ docker run -p 80:80 suckowbiz/fritzbox-cloudflare-dnsupdater
    ```

1. Edit AVM FRITZ!Box: `fritzbox > Internet > Freigaben > DynDNS (Benutzerdefiniert)`:

    - **DynDNS-Anbieter**: `Benutzerdefiniert`
    - **Update-URL**: `http://TODO1/update?token=<pass>&ip=<ipaddr>&zone_id=TODO2`  
      
      Where:
      
      - `<ipaddr` automatically replaced with the new IP
      - `<pass>` automatically replaced with value of `Kennwort`
      - `TODO1` address of the server that runs `fritzbox-cloudflare-dnsupdater`
      - `TODO2` Zone ID of your Site managed at Cloudflare. To support many Cloudflare sites the `zone_id` field can
       be repeated. The Zone ID is present on the Cloudflare Dashboard: `<Site> > Overview > API/Zone ID`.

         The Zone ID can also be fetched from the Cloudflare API with:

         ```bash
         # insert your Cloudflare login email and Cloudflare `Global API Key` (copy API Key from: https://dash.cloudflare.com/profile/api-tokens)
         curl -X GET "https://api.cloudflare.com/client/v4/zones" \
           -H "Content-Type: application/json" \
           -H "X-Auth-Email: TODO" \
           -H "X-Auth-Key: TODO"
         ```
    - **Domainname**: One of the domain names of the Cloudflare site.
    - **Benutzername**: This field must not be empty to satisfy AVM FRITZ!Box form submission. The value is not used by
     fritzbox-cloudflare-dnsupdater
    - **Kennwort**: The Cloudflare `API Token` with permission to edit **Zone DNS**. Generate it here: [https://dash.cloudflare.com/profile/api-tokens](https://dash.cloudflare.com/profile/api-tokens)

## FRITZ!Box limitations addressed

`fritzbox-cloudflare-dnsupdater` solves the following issues with AVM FRITZ!Box:

- Wrapping the **PUT** call required to update a DNS type "A" record at Cloudflare into a **GET**
- Editing multiple DNS address records at Cloudflare using a single call (Cloudflare API requires multiple
 calls of )
- Updating **n** DNS records with a single request (AVM FRITZ!Box can trigger only a single GET; Cloudflare API
 supports no batch calls)

## Resources

- Cloudflare API: [Update DNS Record](https://api.cloudflare.com/#dns-records-for-a-zone-update-dns-record)