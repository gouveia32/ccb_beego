-- --------------------------------------------------------
-- Servidor:                     127.0.0.1
-- Versão do servidor:           5.7.40 - MySQL Community Server (GPL)
-- OS do Servidor:               Win64
-- HeidiSQL Versão:              11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Copiando estrutura do banco de dados para ccb_beego
CREATE DATABASE IF NOT EXISTS `ccb_beego` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `ccb_beego`;

-- Copiando estrutura para tabela ccb_beego.cliente
CREATE TABLE IF NOT EXISTS `cliente` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `nome` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `contato_funcao` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `contato_nome` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `cgc_cpf` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `razao_social` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `inscr_estadual` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `endereco` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `cidade` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `uf` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `cep` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `telefone_1` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `telefone_2` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `telefone_3` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `obs` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `criado_em` datetime DEFAULT NULL,
  `alterado_em` datetime DEFAULT NULL,
  `estado` int(11) NOT NULL DEFAULT '0',
  `preco_base` double DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.cliente: ~28 rows (aproximadamente)
/*!40000 ALTER TABLE `cliente` DISABLE KEYS */;
INSERT INTO `cliente` (`id`, `nome`, `contato_funcao`, `contato_nome`, `cgc_cpf`, `razao_social`, `inscr_estadual`, `endereco`, `cidade`, `uf`, `cep`, `telefone_1`, `telefone_2`, `telefone_3`, `email`, `obs`, `criado_em`, `alterado_em`, `estado`, `preco_base`) VALUES
	(1, 'José Alves de Gouveia', '', '', '', '', '', 'Rua Dr. José Luciano Siqueira, 158 - Suissa', 'Aracaju', 'SE', '49050566', '98801-5232', '', '', 'gouveia32@gmail.com', 'Teste de obs...', '2023-05-08 03:16:51', '2023-05-08 04:33:21', 0, 10.5),
	(2, 'GLOBO DISTRIBUIDORA BROKER LOG. LTDA', '', ' ', '13006150/0001-54', '', '', 'RUA O, 17, DIA', '', '--', '', '', '', '', '', '', NULL, '2023-05-08 04:41:20', 0, 0),
	(3, 'Virgínia [ Iran]', '', ' ', '', '', '', '', '', '--', '', '', '', '', '', '', NULL, '2023-05-08 04:41:41', 0, 0),
	(4, 'SENAC', 'Proprietáro', 'Carlos', '0365.4618/0001-63', 'SERVICO NACIONAL DE APRENDIZAGEM COMERCIAL - SENAC', '', 'Av. Ivo do Prado, B.São José,564', 'Aracaju', 'SE', '     -', '3217-6028', '', '', '', '', NULL, '2023-05-08 04:41:47', 0, 0),
	(5, '3R Confecções Ltda', '', 'Rivaldo/Graça ', '', '', '', 'Praça Clodoaldo Alencar,s/n sala3 e 4 Leite Neto', '', '--', '', '231-1284', '', '', '', '', NULL, '2023-05-08 04:41:53', 0, 0),
	(6, 'A Suprema Com. e Ind Ltda', '', 'Regina ', '', '', '', 'Rua Laranjeiras, 58 - Centro', 'Aracaju', 'SE', '', '79 231-5241', '', '', '', '', NULL, '2023-05-10 04:07:47', 0, 0),
	(8, 'AMANDA SIQUEIRA DE LIMA', '', 'FÁBIO ', '05794150/0001-38', '', '', 'RUA LUIZ CHARGAS , 44A ATALAIA', 'ARACAJU', 'SE', '', '', '', '', '', '', NULL, '2023-05-08 04:43:25', 0, 0),
	(9, 'ANARTE', '', 'ANA ', '13.131.222/0001-95', '', '', 'Av. Hermes Fontes, 1324', '', '--', '     -', '99648-8000', '00248-1918', '', '', '', NULL, '2023-05-10 01:22:22', 0, 0),
	(10, 'Academia Maynara', 'Gerente', 'João Batista ', '', '', '', '', '', '--', '', '251-3034', '', '', '', '', NULL, '2023-05-10 02:48:04', 0, 0),
	(11, 'ELZA IMPERATRIZ BATALHA DE GOIS (Passo a Passo)', '', 'Elzinha ', '32.711.913/0001-02', '', '', 'Av. Jorge Amado, 1055 Galeria Espaço 1055,  - Sl 404, B. Jard', 'Aracaju', 'SE', '     -', '32461-500', '99817-299', '     -', 'academiapassoapasso@hotmail.com', '', NULL, '2023-05-09 05:57:43', 0, 0),
	(12, 'Academia Winner', '', 'Bárbara Salomé Costa Pinto', '', NULL, '', 'Av. Mário Jorge Menezes Vieira, 2266 - Corôa do Meio', 'Aracaju', 'Se', '', '255-2562', '9939-9099', NULL, '', '', NULL, NULL, 1, 0),
	(13, 'Acadêmia Galpão', '', 'Marcos ', '', '', '', 'Rua- Euclides Paes Mendonça, 66 B. 13 de Julho', 'Aracaju', '--', '', '', '', '', '', '', NULL, '2023-05-08 04:42:40', 0, 0),
	(14, 'Acàcia Da Ligue taxi', '99540420', ' ', '', '', '', '', '', '--', '', '', '', '', '', '', NULL, '2023-05-08 04:42:48', 0, 0),
	(15, 'Agropecuária Aquidabã Ltda', '', ' ', '03912192/0001-09', NULL, '127101198', 'Praça da Bandeira Centro,8', 'Aquidabã', 'SE', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(16, 'Agropecuária Vetri Campos Ltda', '', ' ', '32835084/0001-17', NULL, '127681609-7', 'Rua Arquibaldo da Silveira ,36', 'Propriá', 'SE', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(17, 'Algo Mais Confecções', '', '', '17.449.050/0001-53', 'AGPX INDUSTRIA & COMERCIO EIRELI EPP', '27110162-8', 'Rua Engenheiro Marcondes Ferraz, 33 ( ao lado da Energisa)', 'Aracaju', 'SE', '     -', '3217-1726', '', '', 'financeiro@algomaisse.com.br', 'OBS  ENTREGA DE PEÇAS EM PRODUÇÃO\r\n\r\nBORDADO COM 6.000 PONTOS , 5 CORES\r\n\r\nPRODUZIMOS  MAIS OU MENOS 800 PEÇAS POR SEMANA ! \r\n\r\n algomais@infonet.com.br \r\n', NULL, NULL, 1, 0),
	(18, 'Almerinda Jardins Artigos Vestuários Ltda', '', 'Magna ', '093752430001-89', '', '', 'Av. Jorge Amado / 1519 (Bairro Jardins)', 'Aracaju', 'SE', '49025-300', '     -', '     -', '     -', 'almerindajardins@hotmail.com', 'almerindajardins@hotmail.com\r\nReajuste dia 01/12/2017( entre 10 a 9%)', NULL, '2023-05-08 04:42:32', 0, 0),
	(19, 'Alvair Santana Teodoro', '', 'Alvair ', '189696665-91', NULL, '', 'Travesa Ana Célia    No  96    Bairro Bugio', 'Aracaju', 'SE', '49090010', '3252-3038', '', '', '', 'michelcircuitos@yahoo.com.br', NULL, NULL, 1, 0),
	(21, 'Alves Barreto Comercio e Construções Ltda', '', 'Antônio Carlos ', '13004833/0001-15', NULL, '27063274-3', 'Rua Manoel de Oliveira Martins /155A', 'Aracaju', 'SE', '', '2107-2623/2600', '', '', '', 'antoniocarlos@alvesbarreto.com.br', NULL, NULL, 1, 0),
	(22, 'Alves Barreto Comércio e Construções Ltda', '', ' ', '13.004.833/0001-72', NULL, '27.063.247-3', 'Rua Manoel de Oliveira Martins, 155', '', '', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(23, 'Amv Industria Comércio e Serviços de Confecções Ltda', '', ' ', '07646422/0001-88', NULL, '', 'Av. Antônio Fagundes/ Numero 150', 'Aracaju', 'SE', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(24, 'Ana Paula', '', 'Ana Paula ', '', NULL, '', 'Rua dr. José Luciano siqueira, 206 - Antares Apto 103', 'Aracaju', 'Se', '', '224-4524', '', NULL, '', '', NULL, NULL, 1, 0),
	(25, 'Andre  Amigo de monica da Dal', '', ' ', '', NULL, '', '', '', '', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(26, 'Angel Indústria Comércio de confecções Ltda.', 'Gerente', 'Angélica ', '', NULL, '', 'Av Carlos Burlamarque, 226 - Centro', 'Aracaju', 'Se', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(27, 'Antonio  do KM', '', ' ', '', NULL, '', '', '', '', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(28, 'Apple Confecções', '', 'Paulo ', '', '', '', 'Rua João Avila Neto ,54', 'ARACAJU', 'SE', '     -', '23101-54', '99981-7294', '     -', '', 'Próximo a Coca-Cola e ao Extra', NULL, '2023-05-08 05:11:41', 0, 0),
	(29, 'Aracaju Plei Antonio', '', ' ', '', NULL, '', '', '', '', '', '', '', NULL, '', '', NULL, NULL, 1, 0),
	(30, 'Aragão Alfaiataria', '', ' ', '', '', '', 'Rua Campo do Brito / 1362', '', '--', '', '', '', '', '', '\r\naragaoalfaiate@gmail.com', NULL, '2023-05-09 06:01:57', 0, 0);
/*!40000 ALTER TABLE `cliente` ENABLE KEYS */;

-- Copiando estrutura para tabela ccb_beego.sys_backend_user
CREATE TABLE IF NOT EXISTS `sys_backend_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `real_name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `user_name` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `user_pwd` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_super` tinyint(1) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `mobile` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `avatar` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.sys_backend_user: ~2 rows (aproximadamente)
/*!40000 ALTER TABLE `sys_backend_user` DISABLE KEYS */;
INSERT INTO `sys_backend_user` (`id`, `real_name`, `user_name`, `user_pwd`, `is_super`, `status`, `mobile`, `email`, `avatar`) VALUES
	(1, 'Administrador', 'admin', 'e10adc3949ba59abbe56e057f20f883e', 1, 1, '32115232', 'gouveia32@gmail.com', '/static/upload/a290X290.jpg'),
	(2, 'José Gouveia', 'gouveia', 'e10adc3949ba59abbe56e057f20f883e', 0, 1, '79988015232', 'gouveia32@gmail.com', '');
/*!40000 ALTER TABLE `sys_backend_user` ENABLE KEYS */;

-- Copiando estrutura para tabela ccb_beego.sys_logintrace
CREATE TABLE IF NOT EXISTS `sys_logintrace` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `remoteAddr` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `loginTime` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.sys_logintrace: ~0 rows (aproximadamente)
/*!40000 ALTER TABLE `sys_logintrace` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_logintrace` ENABLE KEYS */;

-- Copiando estrutura para tabela ccb_beego.sys_resource
CREATE TABLE IF NOT EXISTS `sys_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `parent_id` int(11) DEFAULT NULL,
  `rtype` int(11) NOT NULL DEFAULT '0',
  `seq` int(11) NOT NULL DEFAULT '0',
  `icon` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `url_for` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=98 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.sys_resource: ~18 rows (aproximadamente)
/*!40000 ALTER TABLE `sys_resource` DISABLE KEYS */;
INSERT INTO `sys_resource` (`id`, `name`, `parent_id`, `rtype`, `seq`, `icon`, `url_for`) VALUES
	(1, 'Sistema', NULL, 1, 30, 'fa fa-gears', ''),
	(2, 'Administração de sistema', 8, 1, 120, 'fa fa-gears', ''),
	(3, 'Instrumento 1', 75, 1, 1, 'fa fa-wpforms', 'HomeController.Index'),
	(5, 'editar', 66, 2, 100, 'fa fa-pencil', 'BackendConfController.Edit'),
	(6, 'excluir', 66, 2, 100, 'fa fa-trash', 'BackendConfController.Delete'),
	(7, 'Ferramenta do sistema', 8, 1, 100, 'fa fa-magnet', ''),
	(8, 'Websocket test', 72, 1, 100, 'fa fa-skyatlas', 'WebsocketWidgetController.Index'),
	(10, 'Medidor 2', 75, 1, 2, 'fa fa-tv', 'HomeController.Index2'),
	(12, 'editar', 84, 2, 100, 'fa fa-pencil', 'SystemValController.Edit'),
	(13, 'excluir', 84, 2, 101, 'fa fa-trash', 'SystemValController.Delete'),
	(14, 'Cadastro', NULL, 1, 1, 'fa fa-address-book-o', ''),
	(15, 'Cliente', 14, 1, 100, 'fa fa-ship', 'ClienteController.Index'),
	(92, 'Ger. dos Recursos', 1, 1, 100, 'fa fa-gear', 'ResourceController.Index'),
	(93, 'Usuários', 1, 1, 100, 'fa fa-reddit', 'BackendUserController.Index'),
	(94, 'Funções', 1, 1, 100, 'fa fa-eyedropper', 'RoleController.Index'),
	(95, 'Alterar Cliente', 15, 2, 100, 'fa fa-user-secret', 'ClienteController.Edit'),
	(96, 'Excluir Cliente', 15, 2, 100, 'fa fa-user-times', 'ClienteController.Delete'),
	(97, 'Alocar Funções', 94, 2, 100, 'fa fa-random', 'RoleController.Allocate');
/*!40000 ALTER TABLE `sys_resource` ENABLE KEYS */;

-- Copiando estrutura para tabela ccb_beego.sys_role
CREATE TABLE IF NOT EXISTS `sys_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `seq` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.sys_role: ~2 rows (aproximadamente)
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` (`id`, `name`, `seq`) VALUES
	(1, 'Atendente', 1),
	(2, 'Gerente', 2);
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;

-- Copiando estrutura para tabela ccb_beego.sys_role_backenduser_rel
CREATE TABLE IF NOT EXISTS `sys_role_backenduser_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `backend_user_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.sys_role_backenduser_rel: ~1 rows (aproximadamente)
/*!40000 ALTER TABLE `sys_role_backenduser_rel` DISABLE KEYS */;
INSERT INTO `sys_role_backenduser_rel` (`id`, `role_id`, `backend_user_id`, `created`) VALUES
	(1, 2, 2, '2023-05-10 01:39:06');
/*!40000 ALTER TABLE `sys_role_backenduser_rel` ENABLE KEYS */;

-- Copiando estrutura para tabela ccb_beego.sys_role_resource_rel
CREATE TABLE IF NOT EXISTS `sys_role_resource_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `resource_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Copiando dados para a tabela ccb_beego.sys_role_resource_rel: ~9 rows (aproximadamente)
/*!40000 ALTER TABLE `sys_role_resource_rel` DISABLE KEYS */;
INSERT INTO `sys_role_resource_rel` (`id`, `role_id`, `resource_id`, `created`) VALUES
	(7, 1, 14, '2023-05-10 02:01:05'),
	(8, 1, 15, '2023-05-10 02:01:05'),
	(18, 2, 14, '2023-05-10 02:14:02'),
	(19, 2, 15, '2023-05-10 02:14:02'),
	(20, 2, 95, '2023-05-10 02:14:02'),
	(21, 2, 96, '2023-05-10 02:14:02'),
	(22, 2, 1, '2023-05-10 02:14:02'),
	(23, 2, 94, '2023-05-10 02:14:02'),
	(24, 2, 97, '2023-05-10 02:14:02');
/*!40000 ALTER TABLE `sys_role_resource_rel` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
