---
page_title: "Resource: auth0_encryption_key_manager"
description: |-
  Resource to allow the rekeying of your tenant master key.
---

# Resource: auth0_encryption_key_manager

Resource to allow the rekeying of your tenant master key.

## Example Usage

```terraform
resource "auth0_encryption_key_manager" "my_encryption_key_manager_initial" {
  key_rotation_id = "da9f2f3b-1c7e-4245-8982-9a25da8407c4"
}

resource "auth0_encryption_key_manager" "my_encryption_key_manager_rekey" {
  key_rotation_id = "68feba2c-7768-40f3-9d71-4b91e0233abf"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `key_rotation_id` (String) If this value is changed, the encryption keys will be rotated. A UUID is recommended for the `key_rotation_id`.

### Read-Only

- `encryption_keys` (List of Object) All encryption keys. (see [below for nested schema](#nestedatt--encryption_keys))
- `id` (String) The ID of this resource.

<a id="nestedatt--encryption_keys"></a>
### Nested Schema for `encryption_keys`

Read-Only:

- `created_at` (String)
- `key_id` (String)
- `parent_key_id` (String)
- `state` (String)
- `type` (String)
- `updated_at` (String)

## Import

Import is supported using the following syntax:

```shell
# As this is not a resource identifiable by an ID within the Auth0 Management API,
# auth0_encryption_key_manager can be imported using a random string.
#
# We recommend [Version 4 UUID](https://www.uuidgenerator.net/version4)
#
# Example:
terraform import auth0_encryption_key_manager.my_key_manager "6f0519ad-ea35-44a3-9b0e-ac9c631612c2"
```