-- Clear all data except users
DELETE FROM audit_entries;
DELETE FROM flag_environments;
DELETE FROM flags;
DELETE FROM environments;
DELETE FROM projects;

-- Projects
INSERT INTO projects (id, name, key, description, created_by, created_at, updated_at) VALUES
  (1, 'Web App',     'web-app',     'Customer-facing web application',          2, datetime('now', '-30 days'), datetime('now', '-2 days')),
  (2, 'Mobile App',  'mobile-app',  'iOS and Android mobile client',            2, datetime('now', '-20 days'), datetime('now', '-1 day')),
  (3, 'API Backend', 'api-backend', 'Core REST API and microservices platform', 2, datetime('now', '-10 days'), datetime('now'));

-- Environments (3 per project)
INSERT INTO environments (id, project_id, name, key, description, sdk_key, created_at, updated_at) VALUES
  -- Web App
  (1,  1, 'Development', 'development', 'Local and CI development',          'sdk_webdev000000000000000000000000000000000000000', datetime('now', '-30 days'), datetime('now', '-2 days')),
  (2,  1, 'Staging',     'staging',     'Pre-production validation',         'sdk_webstg000000000000000000000000000000000000000', datetime('now', '-30 days'), datetime('now', '-1 day')),
  (3,  1, 'Production',  'production',  'Live traffic',                      'sdk_webprd000000000000000000000000000000000000000', datetime('now', '-30 days'), datetime('now')),
  -- Mobile App
  (4,  2, 'Development', 'development', 'Local and CI development',          'sdk_mobdev000000000000000000000000000000000000000', datetime('now', '-20 days'), datetime('now', '-1 day')),
  (5,  2, 'Staging',     'staging',     'Pre-production validation',         'sdk_mobstg000000000000000000000000000000000000000', datetime('now', '-20 days'), datetime('now')),
  (6,  2, 'Production',  'production',  'Live traffic',                      'sdk_mobprd000000000000000000000000000000000000000', datetime('now', '-20 days'), datetime('now')),
  -- API Backend
  (7,  3, 'Development', 'development', 'Local and CI development',          'sdk_apidev000000000000000000000000000000000000000', datetime('now', '-10 days'), datetime('now')),
  (8,  3, 'Staging',     'staging',     'Pre-production validation',         'sdk_apistg000000000000000000000000000000000000000', datetime('now', '-10 days'), datetime('now')),
  (9,  3, 'Production',  'production',  'Live traffic',                      'sdk_apiprd000000000000000000000000000000000000000', datetime('now', '-10 days'), datetime('now'));

-- Flags for Web App (project 1)
INSERT INTO flags (id, project_id, key, name, description, flag_type, variations, created_at, updated_at) VALUES
  (1,  1, 'dark-mode',        'Dark Mode',              'Enable dark mode for the UI',                         'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-28 days'), datetime('now', '-2 days')),
  (2,  1, 'new-dashboard',    'New Dashboard',          'Redesigned analytics dashboard',                      'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-25 days'), datetime('now', '-1 day')),
  (3,  1, 'welcome-message',  'Welcome Message',        'Greeting shown on the home page',                     'string',  '[{"name":"default","value":"Welcome back!"},{"name":"promo","value":"Summer sale — 20% off!"}]',                    datetime('now', '-22 days'), datetime('now')),
  (4,  1, 'theme-color',      'Theme Color',            'Primary brand color',                                 'string',  '[{"name":"blue","value":"#2563eb"},{"name":"purple","value":"#7c3aed"},{"name":"green","value":"#16a34a"}]',          datetime('now', '-20 days'), datetime('now')),
  (5,  1, 'items-per-page',   'Items Per Page',         'Default pagination page size',                        'number',  '[{"name":"10","value":10},{"name":"25","value":25},{"name":"50","value":50}]',                                       datetime('now', '-18 days'), datetime('now')),
  (6,  1, 'max-upload-mb',    'Max Upload Size (MB)',   'Maximum file upload size in megabytes',               'number',  '[{"name":"10mb","value":10},{"name":"50mb","value":50},{"name":"100mb","value":100}]',                               datetime('now', '-15 days'), datetime('now')),
  (7,  1, 'maintenance-mode', 'Maintenance Mode',       'Take the site offline for maintenance',               'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-12 days'), datetime('now')),
  (8,  1, 'beta-features',    'Beta Features',          'Expose experimental features to beta users',          'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-10 days'), datetime('now')),
  (9,  1, 'ab-test-variant',  'A/B Test Variant',       'Homepage hero copy variant',                          'string',  '[{"name":"control","value":"control"},{"name":"variant-a","value":"variant-a"},{"name":"variant-b","value":"variant-b"}]', datetime('now', '-8 days'), datetime('now')),
  (10, 1, 'feature-config',   'Feature Config',         'JSON config controlling feature behaviour',           'json',    '[{"name":"default","value":{"sidebar":true,"notifications":true,"analytics":false}},{"name":"full","value":{"sidebar":true,"notifications":true,"analytics":true}}]', datetime('now', '-5 days'), datetime('now'));

-- Flags for Mobile App (project 2)
INSERT INTO flags (id, project_id, key, name, description, flag_type, variations, created_at, updated_at) VALUES
  (11, 2, 'dark-mode',        'Dark Mode',              'Enable dark mode for the app',                        'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-18 days'), datetime('now', '-1 day')),
  (12, 2, 'onboarding-v2',   'Onboarding v2',          'New onboarding flow for fresh installs',              'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-17 days'), datetime('now')),
  (13, 2, 'push-prompt-copy', 'Push Prompt Copy',      'Text shown in the push notification permission prompt','string',  '[{"name":"default","value":"Stay in the loop — enable notifications"},{"name":"short","value":"Enable notifications"}]', datetime('now', '-15 days'), datetime('now')),
  (14, 2, 'app-theme',        'App Theme',              'Visual theme for the app shell',                      'string',  '[{"name":"light","value":"light"},{"name":"dark","value":"dark"},{"name":"system","value":"system"}]',               datetime('now', '-13 days'), datetime('now')),
  (15, 2, 'feed-page-size',   'Feed Page Size',         'Number of items loaded per feed page',                'number',  '[{"name":"20","value":20},{"name":"40","value":40}]',                                                               datetime('now', '-12 days'), datetime('now')),
  (16, 2, 'video-quality',    'Default Video Quality',  'Default streaming video quality in kbps',             'number',  '[{"name":"720p","value":720},{"name":"1080p","value":1080},{"name":"4k","value":2160}]',                            datetime('now', '-10 days'), datetime('now')),
  (17, 2, 'crash-reporting',  'Crash Reporting',        'Send crash reports to Sentry',                        'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-8 days'),  datetime('now')),
  (18, 2, 'beta-features',    'Beta Features',          'Experimental features for opted-in beta users',       'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-6 days'),  datetime('now')),
  (19, 2, 'ab-checkout-flow', 'A/B Checkout Flow',      'Checkout UX experiment',                              'string',  '[{"name":"control","value":"control"},{"name":"simplified","value":"simplified"}]',                                  datetime('now', '-4 days'),  datetime('now')),
  (20, 2, 'remote-config',    'Remote Config',          'JSON config bundle pushed to clients',                'json',    '[{"name":"default","value":{"offlineMode":false,"maxCacheAgeDays":7}},{"name":"aggressive-cache","value":{"offlineMode":true,"maxCacheAgeDays":30}}]', datetime('now', '-2 days'), datetime('now'));

-- Flags for API Backend (project 3)
INSERT INTO flags (id, project_id, key, name, description, flag_type, variations, created_at, updated_at) VALUES
  (21, 3, 'rate-limiting',     'Rate Limiting',          'Enforce per-user API rate limits',                   'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-9 days'),  datetime('now')),
  (22, 3, 'requests-per-min',  'Requests Per Minute',    'Default rate limit cap per user',                    'number',  '[{"name":"60","value":60},{"name":"120","value":120},{"name":"300","value":300}]',                                  datetime('now', '-9 days'),  datetime('now')),
  (23, 3, 'db-pool-size',      'DB Pool Size',           'Max number of database connections in the pool',     'number',  '[{"name":"10","value":10},{"name":"25","value":25},{"name":"50","value":50}]',                                      datetime('now', '-8 days'),  datetime('now')),
  (24, 3, 'auth-provider',     'Auth Provider',          'Which auth provider to use for JWT issuance',        'string',  '[{"name":"internal","value":"internal"},{"name":"auth0","value":"auth0"}]',                                          datetime('now', '-7 days'),  datetime('now')),
  (25, 3, 'log-level',         'Log Level',              'Runtime log verbosity',                              'string',  '[{"name":"info","value":"info"},{"name":"debug","value":"debug"},{"name":"warn","value":"warn"}]',                    datetime('now', '-7 days'),  datetime('now')),
  (26, 3, 'maintenance-mode',  'Maintenance Mode',       'Return 503 for all API requests',                    'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-6 days'),  datetime('now')),
  (27, 3, 'new-query-engine',  'New Query Engine',       'Use the rewritten query planner',                    'boolean', '[{"name":"on","value":true},{"name":"off","value":false}]',                                                         datetime('now', '-5 days'),  datetime('now')),
  (28, 3, 'cache-ttl-seconds', 'Cache TTL (seconds)',    'Response cache time-to-live in seconds',             'number',  '[{"name":"60s","value":60},{"name":"300s","value":300},{"name":"3600s","value":3600}]',                             datetime('now', '-4 days'),  datetime('now')),
  (29, 3, 'webhook-version',   'Webhook Payload Version','Outbound webhook schema version',                    'string',  '[{"name":"v1","value":"v1"},{"name":"v2","value":"v2"}]',                                                            datetime('now', '-2 days'),  datetime('now')),
  (30, 3, 'service-config',    'Service Config',         'JSON config applied to the running service',         'json',    '[{"name":"default","value":{"tracing":false,"metricsPort":9090}},{"name":"observability","value":{"tracing":true,"metricsPort":9090}}]', datetime('now', '-1 day'), datetime('now'));

-- Flag environments: every flag × every environment in its project
-- Web App flags (1–10) × envs 1,2,3
INSERT INTO flag_environments (flag_id, environment_id, enabled, default_variation) VALUES
  -- dark-mode: on everywhere
  (1, 1, 1, 0), (1, 2, 1, 0), (1, 3, 0, 1),
  -- new-dashboard: dev+staging on, prod off
  (2, 1, 1, 0), (2, 2, 1, 0), (2, 3, 0, 1),
  -- welcome-message: promo in dev, default elsewhere
  (3, 1, 1, 1), (3, 2, 1, 0), (3, 3, 1, 0),
  -- theme-color: purple in dev, blue in staging+prod
  (4, 1, 1, 1), (4, 2, 1, 0), (4, 3, 1, 0),
  -- items-per-page: 50 in dev, 25 in staging, 10 in prod
  (5, 1, 1, 2), (5, 2, 1, 1), (5, 3, 1, 0),
  -- max-upload-mb: 100 in dev, 50 in staging, 10 in prod
  (6, 1, 1, 2), (6, 2, 1, 1), (6, 3, 1, 0),
  -- maintenance-mode: off everywhere
  (7, 1, 0, 1), (7, 2, 0, 1), (7, 3, 0, 1),
  -- beta-features: on in dev, off elsewhere
  (8, 1, 1, 0), (8, 2, 0, 1), (8, 3, 0, 1),
  -- ab-test-variant: variant-a in dev, control elsewhere
  (9, 1, 1, 1), (9, 2, 1, 0), (9, 3, 1, 0),
  -- feature-config: full in dev, default elsewhere
  (10, 1, 1, 1), (10, 2, 1, 0), (10, 3, 1, 0);

-- Mobile App flags (11–20) × envs 4,5,6
INSERT INTO flag_environments (flag_id, environment_id, enabled, default_variation) VALUES
  (11, 4, 1, 0), (11, 5, 1, 0), (11, 6, 0, 1),
  (12, 4, 1, 0), (12, 5, 1, 0), (12, 6, 0, 1),
  (13, 4, 1, 0), (13, 5, 1, 0), (13, 6, 1, 0),
  (14, 4, 1, 2), (14, 5, 1, 0), (14, 6, 1, 0),
  (15, 4, 1, 1), (15, 5, 1, 1), (15, 6, 1, 0),
  (16, 4, 1, 1), (16, 5, 1, 1), (16, 6, 1, 0),
  (17, 4, 1, 0), (17, 5, 1, 0), (17, 6, 1, 0),
  (18, 4, 1, 0), (18, 5, 0, 1), (18, 6, 0, 1),
  (19, 4, 1, 1), (19, 5, 1, 0), (19, 6, 1, 0),
  (20, 4, 1, 1), (20, 5, 1, 0), (20, 6, 1, 0);

-- API Backend flags (21–30) × envs 7,8,9
INSERT INTO flag_environments (flag_id, environment_id, enabled, default_variation) VALUES
  (21, 7, 0, 1), (21, 8, 1, 0), (21, 9, 1, 0),
  (22, 7, 1, 2), (22, 8, 1, 1), (22, 9, 1, 0),
  (23, 7, 1, 2), (23, 8, 1, 1), (23, 9, 1, 0),
  (24, 7, 1, 0), (24, 8, 1, 0), (24, 9, 1, 0),
  (25, 7, 1, 1), (25, 8, 1, 0), (25, 9, 1, 0),
  (26, 7, 0, 1), (26, 8, 0, 1), (26, 9, 0, 1),
  (27, 7, 1, 0), (27, 8, 1, 0), (27, 9, 0, 1),
  (28, 7, 1, 2), (28, 8, 1, 1), (28, 9, 1, 0),
  (29, 7, 1, 0), (29, 8, 1, 0), (29, 9, 1, 0),
  (30, 7, 1, 1), (30, 8, 1, 0), (30, 9, 1, 0);
