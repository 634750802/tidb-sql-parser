CREATE TABLE `access_logs` (
                               `id` bigint(20) NOT NULL /*T![auto_rand] AUTO_RANDOM(5) */,
                               `remote_addr` varchar(128) NOT NULL DEFAULT '',
                               `origin` varchar(128) NOT NULL DEFAULT '',
                               `status_code` int(11) NOT NULL DEFAULT '0',
                               `request_path` varchar(256) NOT NULL DEFAULT '',
                               `request_params` json DEFAULT NULL,
                               `requested_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /*T![auto_rand_base] AUTO_RANDOM_BASE=2070002 */;

CREATE TABLE `ar_internal_metadata` (
                                        `key` varchar(255) NOT NULL,
                                        `value` varchar(255) DEFAULT NULL,
                                        `created_at` datetime(6) NOT NULL,
                                        `updated_at` datetime(6) NOT NULL,
                                        PRIMARY KEY (`key`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `blacklist_repos` (
    `name` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `blacklist_users` (
                                   `login` varchar(255) NOT NULL,
                                   UNIQUE KEY `blacklist_users_login_uindex` (`login`),
                                   PRIMARY KEY (`login`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `cache` (
                         `cache_key` varchar(512) NOT NULL,
                         `cache_value` json NOT NULL,
                         `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         `expires` int(11) DEFAULT '-1' COMMENT 'cache will expire after n seconds',
                         PRIMARY KEY (`cache_key`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `cached_table_cache` (
                                      `cache_key` varchar(512) NOT NULL,
                                      `cache_value` json NOT NULL,
                                      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                      `expires` int(11) DEFAULT '-1' COMMENT 'cache will expire after n seconds',
                                      PRIMARY KEY (`cache_key`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `cn_orgs` (
                           `id` varchar(255) NOT NULL,
                           `name` varchar(255) DEFAULT NULL,
                           `company` varchar(255) DEFAULT NULL,
                           PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `cn_repos` (
                            `id` varchar(255) NOT NULL,
                            `name` varchar(255) DEFAULT NULL,
                            `company` varchar(255) DEFAULT NULL,
                            PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `collection_items` (
                                    `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                    `collection_id` bigint(20) DEFAULT NULL,
                                    `repo_name` varchar(255) NOT NULL,
                                    `repo_id` bigint(20) NOT NULL,
                                    `last_month_rank` int(11) DEFAULT NULL,
                                    `last_2nd_month_rank` int(11) DEFAULT NULL,
                                    PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                                    KEY `index_collection_items_on_collection_id` (`collection_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=720018;

CREATE TABLE `collections` (
                               `id` bigint(20) NOT NULL AUTO_INCREMENT,
                               `name` varchar(255) NOT NULL,
                               `public` tinyint(1) DEFAULT '1',
                               PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                               UNIQUE KEY `index_collections_on_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=10064;

CREATE TABLE `coss_invest` (
                               `id` int(11) NOT NULL AUTO_INCREMENT,
                               `company` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `repository` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `stage` enum('Seed','A','B','C','D','E','F','G','Growth') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `round_size` decimal(10,2) DEFAULT NULL,
                               `year` year(4) DEFAULT NULL,
                               `month` int(11) DEFAULT NULL,
                               `lead_investor` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `has_repo` tinyint(1) DEFAULT NULL,
                               `has_github` tinyint(1) DEFAULT NULL,
                               `github_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `date` datetime DEFAULT NULL,
                               PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=90003;

CREATE TABLE `csdn_events` (
                               `id` bigint(20) DEFAULT NULL,
                               `type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `created_at` datetime DEFAULT NULL,
                               `repo_id` bigint(20) DEFAULT NULL,
                               `repo_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `actor_id` bigint(20) DEFAULT NULL,
                               `actor_login` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `actor_location` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `language` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `additions` bigint(20) DEFAULT NULL,
                               `deletions` bigint(20) DEFAULT NULL,
                               `action` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `number` int(11) DEFAULT NULL,
                               `commit_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `comment_id` bigint(20) DEFAULT NULL,
                               `org_login` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `org_id` bigint(20) DEFAULT NULL,
                               `state` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `closed_at` datetime DEFAULT NULL,
                               `comments` int(11) DEFAULT NULL,
                               `pr_merged_at` datetime DEFAULT NULL,
                               `pr_merged` tinyint(1) DEFAULT NULL,
                               `pr_changed_files` int(11) DEFAULT NULL,
                               `pr_review_comments` int(11) DEFAULT NULL,
                               `pr_or_issue_id` bigint(20) DEFAULT NULL,
                               `event_day` date DEFAULT NULL,
                               `event_month` date DEFAULT NULL,
                               `author_association` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `event_year` int(11) DEFAULT NULL,
                               `push_size` int(11) DEFAULT NULL,
                               `push_distinct_size` int(11) DEFAULT NULL,
                               `creator_user_login` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                               `creator_user_id` bigint(20) DEFAULT NULL,
                               `pr_or_issue_created_at` datetime DEFAULT NULL,
                               KEY `index_github_events_on_id` (`id`) /*!80000 INVISIBLE */,
                               KEY `index_github_events_on_actor_login` (`actor_login`),
                               KEY `index_github_events_on_created_at` (`created_at`),
                               KEY `index_github_events_on_repo_name` (`repo_name`),
                               KEY `index_github_events_on_repo_id_type_action_month_actor_login` (`repo_id`,`type`,`action`,`event_month`,`actor_login`),
                               KEY `index_ge_on_repo_id_type_action_pr_merged_created_at_add_del` (`repo_id`,`type`,`action`,`pr_merged`,`created_at`,`additions`,`deletions`),
                               KEY `index_ge_on_creator_id_type_action_merged_created_at_add_del` (`creator_user_id`,`type`,`action`,`pr_merged`,`created_at`,`additions`,`deletions`),
                               KEY `index_ge_on_actor_id_type_action_created_at_repo_id_commits` (`actor_id`,`type`,`action`,`created_at`,`repo_id`,`push_distinct_size`),
                               KEY `index_ge_on_org_id_type_action_pr_merged_created_at_add_del` (`org_id`,`type`,`action`,`pr_merged`,`created_at`,`additions`,`deletions`),
                               KEY `index_ge_on_repo_id_type_action_created_at_number_pdsize_psize` (`repo_id`,`type`,`action`,`created_at`,`number`,`push_distinct_size`,`push_size`),
                               KEY `index_ge_on_org_id_type_action_created_at_number_pdsize_psize` (`org_id`,`type`,`action`,`created_at`,`number`,`push_distinct_size`,`push_size`),
                               KEY `index_github_events_on_org_id_type_action_month_actor_login` (`org_id`,`type`,`action`,`event_month`,`actor_login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY LIST COLUMNS(`type`)
(PARTITION `push_event` VALUES IN ('PushEvent'),
 PARTITION `create_event` VALUES IN ('CreateEvent'),
 PARTITION `pull_request_event` VALUES IN ('PullRequestEvent'),
 PARTITION `watch_event` VALUES IN ('WatchEvent'),
 PARTITION `issue_comment_event` VALUES IN ('IssueCommentEvent'),
 PARTITION `issues_event` VALUES IN ('IssuesEvent'),
 PARTITION `delete_event` VALUES IN ('DeleteEvent'),
 PARTITION `fork_event` VALUES IN ('ForkEvent'),
 PARTITION `pull_request_review_comment_event` VALUES IN ('PullRequestReviewCommentEvent'),
 PARTITION `pull_request_review_event` VALUES IN ('PullRequestReviewEvent'),
 PARTITION `gollum_event` VALUES IN ('GollumEvent'),
 PARTITION `release_event` VALUES IN ('ReleaseEvent'),
 PARTITION `member_event` VALUES IN ('MemberEvent'),
 PARTITION `commit_comment_event` VALUES IN ('CommitCommentEvent'),
 PARTITION `public_event` VALUES IN ('PublicEvent'),
 PARTITION `gist_event` VALUES IN ('GistEvent'),
 PARTITION `follow_event` VALUES IN ('FollowEvent'),
 PARTITION `event` VALUES IN ('Event'),
 PARTITION `download_event` VALUES IN ('DownloadEvent'),
 PARTITION `team_add_event` VALUES IN ('TeamAddEvent'),
 PARTITION `fork_apply_event` VALUES IN ('ForkApplyEvent'));

CREATE TABLE `csdn_repos` (
                              `repo_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                              PRIMARY KEY (`repo_name`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `css_framework_repos` (
                                       `id` varchar(255) NOT NULL,
                                       `name` varchar(255) DEFAULT NULL,
                                       PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `db_repos` (
                            `id` varchar(255) NOT NULL,
                            `name` varchar(255) DEFAULT NULL,
                            PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `event_logs` (
                              `id` bigint(20) NOT NULL AUTO_INCREMENT,
                              `created_at` datetime NOT NULL,
                              PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                              KEY `index_event_logs_on_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=24745776057;

CREATE TABLE `github_events` (
                                 `id` bigint(20) NOT NULL DEFAULT '0',
                                 `type` varchar(29) NOT NULL DEFAULT 'Event',
                                 `created_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
                                 `repo_id` bigint(20) NOT NULL DEFAULT '0',
                                 `repo_name` varchar(140) NOT NULL DEFAULT '',
                                 `actor_id` bigint(20) NOT NULL DEFAULT '0',
                                 `actor_login` varchar(40) NOT NULL DEFAULT '',
                                 `language` varchar(26) NOT NULL DEFAULT '',
                                 `additions` bigint(20) NOT NULL DEFAULT '0',
                                 `deletions` bigint(20) NOT NULL DEFAULT '0',
                                 `action` varchar(11) NOT NULL DEFAULT '',
                                 `number` int(11) NOT NULL DEFAULT '0',
                                 `commit_id` varchar(40) NOT NULL DEFAULT '',
                                 `comment_id` bigint(20) NOT NULL DEFAULT '0',
                                 `org_login` varchar(40) NOT NULL DEFAULT '',
                                 `org_id` bigint(20) NOT NULL DEFAULT '0',
                                 `state` varchar(6) NOT NULL DEFAULT '',
                                 `closed_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
                                 `comments` int(11) NOT NULL DEFAULT '0',
                                 `pr_merged_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
                                 `pr_merged` tinyint(1) NOT NULL DEFAULT '0',
                                 `pr_changed_files` int(11) NOT NULL DEFAULT '0',
                                 `pr_review_comments` int(11) NOT NULL DEFAULT '0',
                                 `pr_or_issue_id` bigint(20) NOT NULL DEFAULT '0',
                                 `event_day` date NOT NULL,
                                 `event_month` date NOT NULL,
                                 `event_year` int(11) NOT NULL,
                                 `push_size` int(11) NOT NULL DEFAULT '0',
                                 `push_distinct_size` int(11) NOT NULL DEFAULT '0',
                                 `creator_user_login` varchar(40) NOT NULL DEFAULT '',
                                 `creator_user_id` bigint(20) NOT NULL DEFAULT '0',
                                 `pr_or_issue_created_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
                                 KEY `index_github_events_on_id` (`id`) /*!80000 INVISIBLE */,
                                 KEY `index_github_events_on_actor_login` (`actor_login`),
                                 KEY `index_github_events_on_created_at` (`created_at`),
                                 KEY `index_github_events_on_repo_name` (`repo_name`),
                                 KEY `index_github_events_on_repo_id_type_action_month_actor_login` (`repo_id`,`type`,`action`,`event_month`,`actor_login`),
                                 KEY `index_ge_on_repo_id_type_action_pr_merged_created_at_add_del` (`repo_id`,`type`,`action`,`pr_merged`,`created_at`,`additions`,`deletions`),
                                 KEY `index_ge_on_creator_id_type_action_merged_created_at_add_del` (`creator_user_id`,`type`,`action`,`pr_merged`,`created_at`,`additions`,`deletions`),
                                 KEY `index_ge_on_actor_id_type_action_created_at_repo_id_commits` (`actor_id`,`type`,`action`,`created_at`,`repo_id`,`push_distinct_size`),
                                 KEY `index_ge_on_org_id_type_action_pr_merged_created_at_add_del` (`org_id`,`type`,`action`,`pr_merged`,`created_at`,`additions`,`deletions`),
                                 KEY `index_ge_on_repo_id_type_action_created_at_number_pdsize_psize` (`repo_id`,`type`,`action`,`created_at`,`number`,`push_distinct_size`,`push_size`),
                                 KEY `index_ge_on_org_id_type_action_created_at_number_pdsize_psize` (`org_id`,`type`,`action`,`created_at`,`number`,`push_distinct_size`,`push_size`),
                                 KEY `index_github_events_on_org_id_type_action_month_actor_login` (`org_id`,`type`,`action`,`event_month`,`actor_login`),
                                 KEY `index_ge_on_repo_id_type_action_created_at_actor_login` (`repo_id`,`type`,`action`,`created_at`,`actor_login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
PARTITION BY LIST COLUMNS(`type`)
(PARTITION `push_event` VALUES IN ('PushEvent'),
 PARTITION `create_event` VALUES IN ('CreateEvent'),
 PARTITION `pull_request_event` VALUES IN ('PullRequestEvent'),
 PARTITION `watch_event` VALUES IN ('WatchEvent'),
 PARTITION `issue_comment_event` VALUES IN ('IssueCommentEvent'),
 PARTITION `issues_event` VALUES IN ('IssuesEvent'),
 PARTITION `delete_event` VALUES IN ('DeleteEvent'),
 PARTITION `fork_event` VALUES IN ('ForkEvent'),
 PARTITION `pull_request_review_comment_event` VALUES IN ('PullRequestReviewCommentEvent'),
 PARTITION `pull_request_review_event` VALUES IN ('PullRequestReviewEvent'),
 PARTITION `gollum_event` VALUES IN ('GollumEvent'),
 PARTITION `release_event` VALUES IN ('ReleaseEvent'),
 PARTITION `member_event` VALUES IN ('MemberEvent'),
 PARTITION `commit_comment_event` VALUES IN ('CommitCommentEvent'),
 PARTITION `public_event` VALUES IN ('PublicEvent'),
 PARTITION `gist_event` VALUES IN ('GistEvent'),
 PARTITION `follow_event` VALUES IN ('FollowEvent'),
 PARTITION `event` VALUES IN ('Event'),
 PARTITION `download_event` VALUES IN ('DownloadEvent'),
 PARTITION `team_add_event` VALUES IN ('TeamAddEvent'),
 PARTITION `fork_apply_event` VALUES IN ('ForkApplyEvent'));

CREATE TABLE `github_repo_languages` (
                                         `repo_id` int(11) NOT NULL,
                                         `language` varchar(32) NOT NULL,
                                         `size` bigint(20) NOT NULL DEFAULT '0',
                                         PRIMARY KEY (`repo_id`,`language`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `github_repo_topics` (
                                      `repo_id` int(11) NOT NULL,
                                      `topic` varchar(50) NOT NULL,
                                      PRIMARY KEY (`repo_id`,`topic`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `github_repos` (
                                `repo_id` int(11) NOT NULL,
                                `repo_name` varchar(150) NOT NULL,
                                `owner_id` int(11) NOT NULL,
                                `owner_login` varchar(255) NOT NULL,
                                `owner_is_org` tinyint(1) NOT NULL,
                                `description` varchar(512) NOT NULL DEFAULT '',
                                `primary_language` varchar(32) NOT NULL DEFAULT '',
                                `license` varchar(32) NOT NULL DEFAULT '',
                                `size` bigint(20) NOT NULL DEFAULT '0',
                                `stars` int(11) NOT NULL DEFAULT '0',
                                `forks` int(11) NOT NULL DEFAULT '0',
                                `parent_repo_id` int(11) DEFAULT NULL,
                                `is_fork` tinyint(1) NOT NULL DEFAULT '0',
                                `is_archived` tinyint(1) NOT NULL DEFAULT '0',
                                `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
                                `latest_released_at` timestamp NULL DEFAULT NULL,
                                `pushed_at` timestamp NULL DEFAULT NULL,
                                `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                `last_event_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                `refreshed_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                PRIMARY KEY (`repo_id`) /*T![clustered_index] CLUSTERED */,
                                KEY `index_gr_on_owner_id` (`owner_id`),
                                KEY `index_gr_on_repo_name` (`repo_name`),
                                KEY `index_gr_on_stars` (`stars`),
                                KEY `index_gr_on_repo_id_repo_name` (`repo_id`,`repo_name`),
                                KEY `index_gr_on_created_at_is_deleted` (`created_at`,`is_deleted`),
                                KEY `index_gr_on_owner_login_owner_id_is_deleted` (`owner_login`,`owner_id`,`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `github_users` (
                                `id` int(11) NOT NULL,
                                `login` varchar(255) NOT NULL,
                                `type` char(3) NOT NULL DEFAULT 'N/A',
                                `is_bot` tinyint(1) NOT NULL DEFAULT '0',
                                `name` varchar(255) NOT NULL DEFAULT '',
                                `email` varchar(255) NOT NULL DEFAULT '',
                                `organization` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
                                `organization_formatted` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
                                `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
                                `country_code` char(3) NOT NULL DEFAULT 'N/A',
                                `region_code` char(3) NOT NULL DEFAULT 'N/A',
                                `state` varchar(255) NOT NULL DEFAULT '',
                                `city` varchar(255) NOT NULL DEFAULT '',
                                `longitude` decimal(11,8) NOT NULL DEFAULT '0',
                                `latitude` decimal(10,8) NOT NULL DEFAULT '0',
                                `public_repos` int(11) NOT NULL DEFAULT '0',
                                `followers` int(11) NOT NULL DEFAULT '0',
                                `followings` int(11) NOT NULL DEFAULT '0',
                                `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
                                `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
                                `refreshed_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                                KEY `index_gu_on_login_is_bot_organization_country_code` (`login`,`is_bot`,`organization_formatted`,`country_code`),
                                KEY `index_gu_on_address` (`address`),
                                KEY `index_gu_on_organization` (`organization`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `import_logs` (
                               `id` bigint(20) NOT NULL AUTO_INCREMENT,
                               `filename` varchar(255) NOT NULL,
                               `local_file` varchar(255) DEFAULT NULL,
                               `start_download_at` datetime DEFAULT NULL,
                               `end_download_at` datetime DEFAULT NULL,
                               `start_import_at` datetime DEFAULT NULL,
                               `end_import_at` datetime DEFAULT NULL,
                               `start_batch_at` datetime DEFAULT NULL,
                               `created_at` datetime(6) NOT NULL,
                               `updated_at` datetime(6) NOT NULL,
                               PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                               KEY `index_import_logs_on_filename` (`filename`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=2130006;

CREATE TABLE `js_framework_repos` (
                                      `id` varchar(255) NOT NULL,
                                      `name` varchar(255) DEFAULT NULL,
                                      PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `location_cache` (
                                  `address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                                  `valid` tinyint(1) NOT NULL DEFAULT '1',
                                  `formatted_address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
                                  `country_code` char(3) COLLATE utf8mb4_bin NOT NULL DEFAULT 'N/A',
                                  `region_code` char(3) COLLATE utf8mb4_bin NOT NULL DEFAULT 'N/A',
                                  `state` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                  `city` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
                                  `longitude` decimal(11,8) DEFAULT NULL,
                                  `latitude` decimal(10,8) DEFAULT NULL,
                                  `provider` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'UNKNOWN',
                                  PRIMARY KEY (`address`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `mv_coss_dev_month` (
                                     `github_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                                     `event_month` date NOT NULL,
                                     `event_num` int(11) DEFAULT NULL,
                                     `star_num` int(11) DEFAULT NULL,
                                     `pr_num` int(11) DEFAULT NULL,
                                     `issue_num` int(11) DEFAULT NULL,
                                     `dev_num` int(11) DEFAULT NULL,
                                     `star_dev_num` int(11) DEFAULT NULL,
                                     `pr_dev_num` int(11) DEFAULT NULL,
                                     `issue_dev_num` int(11) DEFAULT NULL,
                                     PRIMARY KEY (`github_name`,`event_month`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `nocode_repos` (
                                `id` varchar(255) NOT NULL,
                                `name` varchar(255) DEFAULT NULL,
                                PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `osdb_repos` (
                              `id` varchar(255) NOT NULL,
                              `name` varchar(255) DEFAULT NULL,
                              `group_name` varchar(255) DEFAULT NULL,
                              PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `programming_language_repos` (
                                              `id` varchar(255) NOT NULL,
                                              `name` varchar(255) DEFAULT NULL,
                                              PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `schema_migrations` (
                                     `version` varchar(255) NOT NULL,
                                     PRIMARY KEY (`version`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `share_data` (
                              `id` char(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL COMMENT 'Share ID, uuid v4 generated by client',
                              `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Content for <title>',
                              `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Content for <meta name=description>',
                              `keyword` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Content for <meta name=keyword>',
                              `image_url` varchar(255) CHARACTER SET ascii COLLATE ascii_bin NOT NULL COMMENT 'Default image url',
                              `path` varchar(2083) CHARACTER SET ascii COLLATE ascii_bin NOT NULL COMMENT 'Real url path (includes hash and querystring)',
                              `meta` json NOT NULL COMMENT 'Meta info',
                              `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              `session_id` char(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL COMMENT 'Session ID generated by client IP',
                              KEY `share_data_idx_path` (`path`),
                              KEY `share_data_idx_session_id` (`session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `static_site_generator_repos` (
                                               `id` varchar(255) NOT NULL,
                                               `name` varchar(255) DEFAULT NULL,
                                               PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `stats_index_summary` (
                                       `summary_begin_time` timestamp NOT NULL,
                                       `summary_end_time` timestamp NOT NULL,
                                       `table_name` varchar(64) NOT NULL,
                                       `index_name` varchar(64) NOT NULL,
                                       `digest` varchar(64) NOT NULL,
                                       `plan_digest` varchar(64) NOT NULL,
                                       `exec_count` bigint(20) unsigned DEFAULT '0',
                                       UNIQUE KEY `unique_sts_on_begin_end_index_digest` (`summary_begin_time`,`summary_end_time`,`index_name`,`digest`,`plan_digest`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `stats_query_summary` (
                                       `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                       `query_name` varchar(128) NOT NULL,
                                       `digest_text` text NOT NULL,
                                       `executed_at` timestamp NOT NULL,
                                       PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                                       KEY `index_sqs_on_executed_at` (`executed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=870001;

CREATE TABLE `trending_repos` (
                                  `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                  `repo_name` varchar(255) DEFAULT NULL,
                                  `created_at` datetime DEFAULT NULL,
                                  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                                  KEY `index_trending_repos_on_repo_name` (`repo_name`),
                                  KEY `index_trending_repos_on_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=1740006;

CREATE TABLE `users` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `login` varchar(255) NOT NULL,
                         `company` varchar(255) DEFAULT NULL,
                         `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         `type` varchar(255) NOT NULL DEFAULT 'USR',
                         `fake` tinyint(1) NOT NULL DEFAULT '0',
                         `deleted` tinyint(1) NOT NULL DEFAULT '0',
                         `long` decimal(11,8) DEFAULT NULL,
                         `lat` decimal(10,8) DEFAULT NULL,
                         `country_code` char(3) DEFAULT NULL,
                         `state` varchar(255) DEFAULT NULL,
                         `city` varchar(255) DEFAULT NULL,
                         `location` varchar(255) DEFAULT NULL,
                         PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
                         KEY `index_login_on_users` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `web_framework_repos` (
                                       `id` varchar(255) NOT NULL,
                                       `name` varchar(255) DEFAULT NULL,
                                       PRIMARY KEY (`id`) /*T![clustered_index] NONCLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

