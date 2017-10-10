package setting


type Data struct {
	Admin struct {
		DisableRegularOrgCreation bool `yaml:"disable_regular_org_creation"`
	} `yaml:"admin"`
	API struct {
		MaxResponseItems int `yaml:"max_response_items"`
	} `yaml:"api"`
	Attachment struct {
		AllowedTypes string `yaml:"allowed_types"`
		Enabled      bool   `yaml:"enabled"`
		MaxFiles     int    `yaml:"max_files"`
		MaxSize      int    `yaml:"max_size"`
		Path         string `yaml:"path"`
	} `yaml:"attachment"`
	Cache struct {
		Adapter  string `yaml:"adapter"`
		Host     string `yaml:"host"`
		Interval int    `yaml:"interval"`
	} `yaml:"cache"`
	Cron struct {
		CheckRepoStats struct {
			RunAtStart bool   `yaml:"run_at_start"`
			Schedule   string `yaml:"schedule"`
		} `yaml:"check_repo_stats"`
		Enabled            bool `yaml:"enabled"`
		RepoArchiveCleanup struct {
			OlderThan  string `yaml:"older_than"`
			RunAtStart bool   `yaml:"run_at_start"`
			Schedule   string `yaml:"schedule"`
		} `yaml:"repo_archive_cleanup"`
		RepoHealthCheck struct {
			Args     string `yaml:"args"`
			Schedule string `yaml:"schedule"`
			Timeout  string `yaml:"timeout"`
		} `yaml:"repo_health_check"`
		RunAtStart    bool `yaml:"run_at_start"`
		UpdateMirrors struct {
			Schedule string `yaml:"schedule"`
		} `yaml:"update_mirrors"`
	} `yaml:"cron"`
	Database struct {
		DbType  string `yaml:"db_type"`
		Host    string `yaml:"host"`
		Name    string `yaml:"name"`
		Passwd  string `yaml:"passwd"`
		Path    string `yaml:"path"`
		SslMode string `yaml:"ssl_mode"`
		User    string `yaml:"user"`
	} `yaml:"database"`
	Git struct {
		DisableDiffHighlight     bool   `yaml:"disable_diff_highlight"`
		GcArgs                   string `yaml:"gc_args"`
		MaxGitDiffFiles          int    `yaml:"max_git_diff_files"`
		MaxGitDiffLineCharacters int    `yaml:"max_git_diff_line_characters"`
		MaxGitDiffLines          int    `yaml:"max_git_diff_lines"`
		Timeout                  struct {
			Clone   int `yaml:"clone"`
			Gc      int `yaml:"gc"`
			Migrate int `yaml:"migrate"`
			Mirror  int `yaml:"mirror"`
			Pull    int `yaml:"pull"`
		} `yaml:"timeout"`
	} `yaml:"git"`
	Global struct {
		AppName string `yaml:"app_name"`
		RunMode string `yaml:"run_mode"`
		RunUser string `yaml:"run_user"`
	} `yaml:"global"`
	HTTP struct {
		AccessControlAllowOrigin string `yaml:"access_control_allow_origin"`
	} `yaml:"http"`
	I18N struct {
		Datelang struct {
			BgBg string `yaml:"bg-bg"`
			CsCz string `yaml:"cs-cz"`
			DeDe string `yaml:"de-de"`
			EnUs string `yaml:"en-us"`
			EsEs string `yaml:"es-es"`
			FiFi string `yaml:"fi-fi"`
			FrFr string `yaml:"fr-fr"`
			GlEs string `yaml:"gl-es"`
			ItIt string `yaml:"it-it"`
			JaJp string `yaml:"ja-jp"`
			KoKr string `yaml:"ko-kr"`
			LvLv string `yaml:"lv-lv"`
			NlNl string `yaml:"nl-nl"`
			PlPl string `yaml:"pl-pl"`
			PtBr string `yaml:"pt-br"`
			RuRu string `yaml:"ru-ru"`
			SrSp string `yaml:"sr-sp"`
			SvSe string `yaml:"sv-se"`
			TrTr string `yaml:"tr-tr"`
			UkUa string `yaml:"uk-ua"`
			ZhCn string `yaml:"zh-cn"`
			ZhHk string `yaml:"zh-hk"`
			ZhTw string `yaml:"zh-tw"`
		} `yaml:"datelang"`
		Langs string `yaml:"langs"`
		Names string `yaml:"names"`
	} `yaml:"i18n"`
	Log struct {
		BufferLen int `yaml:"buffer_len"`
		Console   struct {
			Level string `yaml:"level"`
		} `yaml:"console"`
		File struct {
			DailyRotate  bool   `yaml:"daily_rotate"`
			Level        string `yaml:"level"`
			LogRotate    bool   `yaml:"log_rotate"`
			MaxDays      int    `yaml:"max_days"`
			MaxLines     int    `yaml:"max_lines"`
			MaxSizeShift int    `yaml:"max_size_shift"`
		} `yaml:"file"`
		Level    string `yaml:"level"`
		Mode     string `yaml:"mode"`
		RootPath string `yaml:"root_path"`
		Slack    struct {
			Level string `yaml:"level"`
			URL   string `yaml:"url"`
		} `yaml:"slack"`
		Xorm struct {
			MaxDays     int  `yaml:"max_days"`
			MaxSize     int  `yaml:"max_size"`
			Rotate      bool `yaml:"rotate"`
			RotateDaily bool `yaml:"rotate_daily"`
		} `yaml:"xorm"`
	} `yaml:"log"`
	Mailer struct {
		CertFile       string `yaml:"cert_file"`
		DisableHelo    string `yaml:"disable_helo"`
		Enabled        bool   `yaml:"enabled"`
		From           string `yaml:"from"`
		HeloHostname   string `yaml:"helo_hostname"`
		Host           string `yaml:"host"`
		KeyFile        string `yaml:"key_file"`
		Passwd         string `yaml:"passwd"`
		SendBufferLen  int    `yaml:"send_buffer_len"`
		SkipVerify     string `yaml:"skip_verify"`
		Subject        string `yaml:"subject"`
		UseCertificate bool   `yaml:"use_certificate"`
		UsePlainText   bool   `yaml:"use_plain_text"`
		User           string `yaml:"user"`
	} `yaml:"mailer"`
	Markdown struct {
		CustomURLSchemes    string `yaml:"custom_url_schemes"`
		EnableHardLineBreak bool   `yaml:"enable_hard_line_break"`
		FileExtensions      string `yaml:"file_extensions"`
	} `yaml:"markdown"`
	Mirror struct {
		DefaultInterval int `yaml:"default_interval"`
	} `yaml:"mirror"`
	Other struct {
		ShowFooterBranding         bool `yaml:"show_footer_branding"`
		ShowFooterTemplateLoadTime bool `yaml:"show_footer_template_load_time"`
		ShowFooterVersion          bool `yaml:"show_footer_version"`
	} `yaml:"other"`
	Picture struct {
		AvatarUploadPath      string `yaml:"avatar_upload_path"`
		DisableGravatar       bool   `yaml:"disable_gravatar"`
		EnableFederatedAvatar bool   `yaml:"enable_federated_avatar"`
		GravatarSource        string `yaml:"gravatar_source"`
	} `yaml:"picture"`
	Release struct {
		Attachment struct {
			AllowedTypes string `yaml:"allowed_types"`
			Enabled      bool   `yaml:"enabled"`
			MaxFiles     int    `yaml:"max_files"`
			MaxSize      int    `yaml:"max_size"`
			Path         string `yaml:"path"`
		} `yaml:"attachment"`
	} `yaml:"release"`
	Repository struct {
		AnsiCharset             string `yaml:"ansi_charset"`
		CommitsFetchConcurrency int    `yaml:"commits_fetch_concurrency"`
		DisableHTTPGit          bool   `yaml:"disable_http_git"`
		Editor                  struct {
			LineWrapExtensions   string `yaml:"line_wrap_extensions"`
			PreviewableFileModes string `yaml:"previewable_file_modes"`
		} `yaml:"editor"`
		EnableLocalPathMigration bool   `yaml:"enable_local_path_migration"`
		EnableRawFileRenderMode  bool   `yaml:"enable_raw_file_render_mode"`
		ForcePrivate             bool   `yaml:"force_private"`
		MaxCreationLimit         int    `yaml:"max_creation_limit"`
		MirrorQueueLength        int    `yaml:"mirror_queue_length"`
		PreferredLicenses        string `yaml:"preferred_licenses"`
		PullRequestQueueLength   int    `yaml:"pull_request_queue_length"`
		Root                     string `yaml:"root"`
		ScriptType               string `yaml:"script_type"`
		Upload                   struct {
			AllowedTypes string `yaml:"allowed_types"`
			Enabled      bool   `yaml:"enabled"`
			FileMaxSize  int    `yaml:"file_max_size"`
			MaxFiles     int    `yaml:"max_files"`
			TempPath     string `yaml:"temp_path"`
		} `yaml:"upload"`
	} `yaml:"repository"`
	Security struct {
		CookieRememberName             string `yaml:"cookie_remember_name"`
		CookieSecure                   bool   `yaml:"cookie_secure"`
		CookieUsername                 string `yaml:"cookie_username"`
		EnableLoginStatusCookie        bool   `yaml:"enable_login_status_cookie"`
		InstallLock                    bool   `yaml:"install_lock"`
		LoginRememberDays              int    `yaml:"login_remember_days"`
		LoginStatusCookieName          string `yaml:"login_status_cookie_name"`
		ReverseProxyAuthenticationUser string `yaml:"reverse_proxy_authentication_user"`
		SecretKey                      string `yaml:"secret_key"`
	} `yaml:"security"`
	Server struct {
		AppDataPath          string `yaml:"app_data_path"`
		CertFile             string `yaml:"cert_file"`
		DisableRouterLog     bool   `yaml:"disable_router_log"`
		DisableSSH           bool   `yaml:"disable_ssh"`
		Domain               string `yaml:"domain"`
		EnableGzip           bool   `yaml:"enable_gzip"`
		HTTPAddr             string `yaml:"http_addr"`
		HTTPPort             int    `yaml:"http_port"`
		KeyFile              string `yaml:"key_file"`
		LandingPage          string `yaml:"landing_page"`
		LocalRootURL         string `yaml:"local_root_url"`
		MinimumKeySizeCheck  bool   `yaml:"minimum_key_size_check"`
		OfflineMode          bool   `yaml:"offline_mode"`
		Protocol             string `yaml:"protocol"`
		RootURL              string `yaml:"root_url"`
		SSHDomain            string `yaml:"ssh_domain"`
		SSHKeyTestPath       string `yaml:"ssh_key_test_path"`
		SSHKeygenPath        string `yaml:"ssh_keygen_path"`
		SSHListenHost        string `yaml:"ssh_listen_host"`
		SSHListenPort        int    `yaml:"ssh_listen_port"`
		SSHPort              int    `yaml:"ssh_port"`
		SSHRootPath          string `yaml:"ssh_root_path"`
		SSHServerCiphers     string `yaml:"ssh_server_ciphers"`
		StartSSHServer       bool   `yaml:"start_ssh_server"`
		StaticRootPath       string `yaml:"static_root_path"`
		TLSMinVersion        string `yaml:"tls_min_version"`
		UnixSocketPermission int    `yaml:"unix_socket_permission"`
	} `yaml:"server"`
	Service struct {
		ActiveCodeLiveMinutes              int  `yaml:"active_code_live_minutes"`
		DisableRegistration                bool `yaml:"disable_registration"`
		EnableCaptcha                      bool `yaml:"enable_captcha"`
		EnableNotifyMail                   bool `yaml:"enable_notify_mail"`
		EnableReverseProxyAuthentication   bool `yaml:"enable_reverse_proxy_authentication"`
		EnableReverseProxyAutoRegistration bool `yaml:"enable_reverse_proxy_auto_registration"`
		RegisterEmailConfirm               bool `yaml:"register_email_confirm"`
		RequireSigninView                  bool `yaml:"require_signin_view"`
		ResetPasswdCodeLiveMinutes         int  `yaml:"reset_passwd_code_live_minutes"`
	} `yaml:"service"`
	Session struct {
		CookieName      string `yaml:"cookie_name"`
		CookieSecure    bool   `yaml:"cookie_secure"`
		CsrfCookieName  string `yaml:"csrf_cookie_name"`
		EnableSetCookie bool   `yaml:"enable_set_cookie"`
		GcIntervalTime  int    `yaml:"gc_interval_time"`
		Provider        string `yaml:"provider"`
		ProviderConfig  string `yaml:"provider_config"`
		SessionLifeTime int    `yaml:"session_life_time"`
	} `yaml:"session"`
	Smartypants struct {
		AngledQuotes bool `yaml:"angled_quotes"`
		Dashes       bool `yaml:"dashes"`
		Enabled      bool `yaml:"enabled"`
		Fractions    bool `yaml:"fractions"`
		LatexDashes  bool `yaml:"latex_dashes"`
	} `yaml:"smartypants"`
	SSH struct {
		MinimumKeySizes struct {
			Dsa     int `yaml:"dsa"`
			Ecdsa   int `yaml:"ecdsa"`
			Ed25519 int `yaml:"ed25519"`
			Rsa     int `yaml:"rsa"`
		} `yaml:"minimum_key_sizes"`
	} `yaml:"ssh"`
	Time struct {
		Format string `yaml:"format"`
	} `yaml:"time"`
	UI struct {
		Admin struct {
			NoticePagingNum int `yaml:"notice_paging_num"`
			OrgPagingNum    int `yaml:"org_paging_num"`
			RepoPagingNum   int `yaml:"repo_paging_num"`
			UserPagingNum   int `yaml:"user_paging_num"`
		} `yaml:"admin"`
		ExplorePagingNum   int    `yaml:"explore_paging_num"`
		FeedMaxCommitNum   int    `yaml:"feed_max_commit_num"`
		IssuePagingNum     int    `yaml:"issue_paging_num"`
		MaxDisplayFileSize int    `yaml:"max_display_file_size"`
		ThemeColorMetaTag  string `yaml:"theme_color_meta_tag"`
		User               struct {
			CommitsPagingNum  int `yaml:"commits_paging_num"`
			NewsFeedPagingNum int `yaml:"news_feed_paging_num"`
			RepoPagingNum     int `yaml:"repo_paging_num"`
		} `yaml:"user"`
	} `yaml:"ui"`
	Webhook struct {
		DeliverTimeout int    `yaml:"deliver_timeout"`
		PagingNum      int    `yaml:"paging_num"`
		QueueLength    int    `yaml:"queue_length"`
		SkipTLSVerify  bool   `yaml:"skip_tls_verify"`
		Types          string `yaml:"types"`
	} `yaml:"webhook"`
}
