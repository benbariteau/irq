CREATE TABLE `deleted_quote` (
  `id` int(11) NOT NULL,
  `text` text NOT NULL,
  `score` int(11) NOT NULL,
  `is_offensive` tinyint(1) NOT NULL DEFAULT '0',
  `is_nishbot` tinyint(1) NOT NULL DEFAULT '0',
  `time_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
