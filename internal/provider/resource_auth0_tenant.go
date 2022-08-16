package provider

import (
	"context"
	"net/http"

	"github.com/auth0/go-auth0/management"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	internalValidation "github.com/auth0/terraform-provider-auth0/internal/validation"
)

func newTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTenant,
		ReadContext:   readTenant,
		UpdateContext: updateTenant,
		DeleteContext: deleteTenant,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "With this resource, you can manage Auth0 tenants, including setting logos and support contact " +
			"information, setting error pages, and configuring default tenant behaviors.",
		Schema: map[string]*schema.Schema{
			"change_password": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Configuration settings for change password page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether to use the custom change password page.",
						},
						"html": {
							Type:     schema.TypeString,
							Required: true,
							Description: "HTML format with supported Liquid syntax. " +
								"Customized content of the change password page.",
						},
					},
				},
			},
			"guardian_mfa_page": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Configuration settings for the Guardian MFA page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether to use the custom Guardian page.",
						},
						"html": {
							Type:     schema.TypeString,
							Required: true,
							Description: "HTML format with supported Liquid syntax. " +
								"Customized content of the Guardian page.",
						},
					},
				},
			},
			"default_audience": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "API Audience to use by default for API Authorization flows. This setting is " +
					"equivalent to appending the audience to every authorization request made to the tenant " +
					"for every application.",
			},
			"default_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Name of the connection to be used for Password Grant exchanges. " +
					"Options include `auth0-adldap`, `ad`, `auth0`, `email`, `sms`, `waad`, and `adfs`.",
			},
			"error_page": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Configuration settings for error pages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html": {
							Type:     schema.TypeString,
							Required: true,
							Description: "HTML format with supported Liquid syntax. " +
								"Customized content of the error page.",
						},
						"show_log_link": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicates whether to show the link to logs as part of the default error page.",
						},
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "URL to redirect to when an error occurs rather than showing the default error page.",
						},
					},
				},
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Friendly name for the tenant.",
			},
			"picture_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "URL of logo to be shown for the tenant. Recommended size is 150px x 150px. " +
					"If no URL is provided, the Auth0 logo will be used.",
			},
			"support_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Support email address for authenticating users.",
			},
			"support_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Support URL for authenticating users.",
			},
			"allowed_logout_urls": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "URLs that Auth0 may redirect to after logout.",
			},
			"sandbox_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Selected sandbox version for the extensibility environment, which allows you to " +
					"use custom scripts to extend parts of Auth0's functionality.",
			},
			"session_lifetime": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatAtLeast(0.01),
				Default:      168,
				Description:  "Number of hours during which a session will stay valid.",
			},
			"idle_session_lifetime": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Default:      72,
				ValidateFunc: validation.FloatAtLeast(0.01),
				Description:  "Number of hours during which a session can be inactive before the user must log in again.",
			},
			"enabled_locales": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Description: "Supported locales for the user interface. The first locale in the list will be " +
					"used to set the default locale.",
			},
			"flags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Configuration settings for tenant flags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_client_connections": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether all current connections should be enabled when a new client is created.",
						},
						"enable_apis_section": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the APIs section is enabled for the tenant.",
						},
						"enable_pipeline2": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether advanced API Authorization scenarios are enabled.",
						},
						"enable_dynamic_client_registration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the tenant allows dynamic client registration.",
						},
						"enable_custom_domain_in_emails": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the tenant allows custom domains in emails.",
						},
						"universal_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the tenant uses Universal Login.",
						},
						"enable_legacy_logs_search_v2": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether to use the older v2 legacy logs search.",
						},
						"disable_clickjack_protection_headers": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether classic Universal Login prompts include additional security headers to prevent clickjacking.",
						},
						"enable_public_signup_user_exists_error": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the public sign up process shows a `user_exists` error if the user already exists.",
						},
						"use_scope_descriptions_for_consent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether to use scope descriptions for consent.",
						},
						"allow_legacy_delegation_grant_types": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the legacy delegation endpoint will be enabled for your account (true) or not available (false).",
						},
						"allow_legacy_ro_grant_types": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether the legacy `auth/ro` endpoint (used with resource owner password and passwordless features) will be enabled for your account (true) or not available (false).",
						},
						"allow_legacy_tokeninfo_endpoint": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If enabled, customers can use Tokeninfo Endpoint, otherwise they can not use it.",
						},
						"enable_legacy_profile": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether ID tokens and the userinfo endpoint includes a complete user profile (true) or only OpenID Connect claims (false).",
						},
						"enable_idtoken_api2": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether ID tokens can be used to authorize some types of requests to API v2 (true) or not (false).",
						},
						"no_disclose_enterprise_connections": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Do not Publish Enterprise Connections Information with IdP domains on the lock configuration file.",
						},
						"disable_management_api_sms_obfuscation": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If true, SMS phone numbers will not be obfuscated in Management API GET calls.",
						},
						"enable_adfs_waad_email_verification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "If enabled, users will be presented with an email verification prompt during their first login when using Azure AD or ADFS connections.",
						},
						"revoke_refresh_token_grant": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Delete underlying grant when a refresh token is revoked via the Authentication API.",
						},
						"dashboard_log_streams_next": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enables beta access to log streaming changes.",
						},
						"dashboard_insights_view": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enables new insights activity page view.",
						},
						"disable_fields_map_fix": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Disables SAML fields map fix for bad mappings with repeated attributes.",
						},
					},
				},
			},
			"universal_login": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Configuration settings for Universal Login.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"colors": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Configuration settings for Universal Login colors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"primary": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Primary button background color in hexadecimal.",
									},
									"page_background": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Background color of login pages in hexadecimal.",
									},
								},
							},
						},
					},
				},
			},
			"default_redirection_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					internalValidation.IsURLWithNoFragment,
					validation.IsURLWithScheme([]string{"https"}),
				),
				Description: "The default absolute redirection URI, must be https and cannot contain a fragment.",
			},
			"session_cookie": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Alters behavior of tenant's session cookie. Contains a single `mode` property.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"persistent",
								"non-persistent",
							}, false),
							Description: "Behavior of tenant session cookie. Accepts either \"persistent\" or \"non-persistent\".",
						},
					},
				},
			},
		},
	}
}

func createTenant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(resource.UniqueId())
	return updateTenant(ctx, d, m)
}

func readTenant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*management.Management)
	tenant, err := api.Tenant.Read()
	if err != nil {
		if mErr, ok := err.(management.Error); ok {
			if mErr.Status() == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	result := multierror.Append(
		d.Set("change_password", flattenTenantChangePassword(tenant.ChangePassword)),
		d.Set("guardian_mfa_page", flattenTenantGuardianMFAPage(tenant.GuardianMFAPage)),
		d.Set("default_audience", tenant.DefaultAudience),
		d.Set("default_directory", tenant.DefaultDirectory),
		d.Set("default_redirection_uri", tenant.DefaultRedirectionURI),
		d.Set("friendly_name", tenant.FriendlyName),
		d.Set("picture_url", tenant.PictureURL),
		d.Set("support_email", tenant.SupportEmail),
		d.Set("support_url", tenant.SupportURL),
		d.Set("allowed_logout_urls", tenant.AllowedLogoutURLs),
		d.Set("session_lifetime", tenant.SessionLifetime),
		d.Set("idle_session_lifetime", tenant.IdleSessionLifetime),
		d.Set("sandbox_version", tenant.SandboxVersion),
		d.Set("enabled_locales", tenant.EnabledLocales),
		d.Set("error_page", flattenTenantErrorPage(tenant.ErrorPage)),
		d.Set("flags", flattenTenantFlags(tenant.Flags)),
		d.Set("universal_login", flattenTenantUniversalLogin(tenant.UniversalLogin)),
		d.Set("session_cookie", flattenTenantSessionCookie(tenant.SessionCookie)),
	)

	return diag.FromErr(result.ErrorOrNil())
}

func updateTenant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenant := expandTenant(d)
	api := m.(*management.Management)
	if err := api.Tenant.Update(tenant); err != nil {
		return diag.FromErr(err)
	}

	return readTenant(ctx, d, m)
}

func deleteTenant(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func expandTenant(d *schema.ResourceData) *management.Tenant {
	tenant := &management.Tenant{
		DefaultAudience:       String(d, "default_audience"),
		DefaultDirectory:      String(d, "default_directory"),
		DefaultRedirectionURI: String(d, "default_redirection_uri"),
		FriendlyName:          String(d, "friendly_name"),
		PictureURL:            String(d, "picture_url"),
		SupportEmail:          String(d, "support_email"),
		SupportURL:            String(d, "support_url"),
		AllowedLogoutURLs:     Slice(d, "allowed_logout_urls"),
		SessionLifetime:       Float64(d, "session_lifetime"),
		SandboxVersion:        String(d, "sandbox_version"),
		IdleSessionLifetime:   Float64(d, "idle_session_lifetime", IsNewResource(), HasChange()),
		EnabledLocales:        List(d, "enabled_locales").List(),
		ChangePassword:        expandTenantChangePassword(d),
		GuardianMFAPage:       expandTenantGuardianMFAPage(d),
		ErrorPage:             expandTenantErrorPage(d),
		Flags:                 expandTenantFlags(d.GetRawConfig().GetAttr("flags")),
		UniversalLogin:        expandTenantUniversalLogin(d),
		SessionCookie:         expandTenantSessionCookie(d),
	}

	return tenant
}