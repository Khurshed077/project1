-- phpMyAdmin SQL Dump
-- version 5.1.0
-- https://www.phpmyadmin.net/
--
-- Хост: 127.0.0.1:3306
-- Время создания: Дек 04 2025 г., 01:35
-- Версия сервера: 5.7.33
-- Версия PHP: 7.1.33

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `golang`
--

-- --------------------------------------------------------

--
-- Структура таблицы `articles`
--

CREATE TABLE `articles` (
  `id` int(11) UNSIGNED NOT NULL,
  `title` varchar(100) NOT NULL,
  `anons` varchar(255) NOT NULL,
  `full_text` text NOT NULL,
  `image` varchar(255) DEFAULT NULL,
  `category_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `articles`
--

INSERT INTO `articles` (`id`, `title`, `anons`, `full_text`, `image`, `category_id`) VALUES
(5, 'В России заблокировали Roblox', 'МОСКВА, 3 дек — РИА Новости. В России заблокировали игровую площадку Roblox из-за распространения материалов с оправданием терроризма, говорится в заявлении РКН.', 'МОСКВА, 3 дек — РИА Новости. В России заблокировали игровую площадку Roblox из-за распространения материалов с оправданием терроризма, говорится в заявлении РКН.\r\n\"Роскомнадзор ограничил доступ к американскому интернет-сервису Roblox в связи с выявленными фактами массового и неоднократного распространения материалов с пропагандой и оправданием экстремистской и террористической деятельности, призывов к совершению противоправных действий насильственного характера и пропаганды ЛГБТ-тематики*\", — сообщили в ведомстве.\r\n\r\n\r\nТам отметили, что площадка Roblox популярна у педофилов, которые знакомятся с несовершеннолетними прямо в чатах игры, а затем переходят в реальную жизнь.', 'uploads/2059467224_0_0_3072_1728_768x0_80_0_0_bf2c69aa6bc27e55e13d515687985f37.jpg.webp', 1),
(6, 'Сколько получают веб-разработчики: зарплаты в России в 2025 году', 'Профессия веб-программиста входит в топ самых востребованных IT-специальностей. Рассказываем, от чего зависит доход такого разработчика в России.', 'Кто такой веб-программист: frontend, backend, fullstack\r\nВеб-программист — это специалист, который создаёт и поддерживает сайты и веб-приложения. Но в это понятие входит несколько направлений работы.\r\nFrontend-разработчик отвечает за внешнюю часть сайта — то, что видит и с чем взаимодействует пользователь. Это интерфейсы, кнопки, анимация, формы. Основные задачи фронтендера — сделать сайт удобным, красивым.\r\nBackend-разработчик занимается «внутренней кухней» сайта — серверной логикой, базами данных, авторизацией пользователей, обработкой платежей.\r\nFullstack-разработчик совмещает оба направления. Он способен работать и с интерфейсом, и с серверной частью. Такого специалиста особенно ценят в небольших компаниях и стартапах.\r\nНачать свой путь в IT можно с позиции фронтенд-разработчика. На курсе Практикума можно освоить фронтенд-разработку с нуля за 10 месяцев. Будет много практики на реальных проектах.', 'uploads/lqip.webp', 2);

-- --------------------------------------------------------

--
-- Структура таблицы `categories`
--

CREATE TABLE `categories` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `categories`
--

INSERT INTO `categories` (`id`, `name`) VALUES
(1, 'Новости\r\n'),
(2, 'Статьи\r\n');

-- --------------------------------------------------------

--
-- Структура таблицы `users`
--

CREATE TABLE `users` (
  `id` int(11) UNSIGNED NOT NULL,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password_hash` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `users`
--

INSERT INTO `users` (`id`, `username`, `email`, `password_hash`) VALUES
(1, 'root', 'khurshedmamadshoev@gmail.com', '$2a$10$K7uoUVY4pfn74RlRb0aX6udaGrn97yAjTvCxtUDCE7IemEEQo7Qtm'),
(2, 'root1', 'khurshedmamadshoev@gmail.com', '$2a$10$Jl5/4OjjULJYRYp0Dduacu/3WxpaySjDIug7PL4HN5fRklQXatPtS'),
(3, 'Хуршед', 'khurshedmamadshoev2@gmail.com', '$2a$10$wuOJ0Dr.POfd3pJe/phFsuPLvy94Tfb2S8o28aYWQuWhpXBkqvyaK');

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `articles`
--
ALTER TABLE `articles`
  ADD PRIMARY KEY (`id`),
  ADD KEY `category_id` (`category_id`);

--
-- Индексы таблицы `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `articles`
--
ALTER TABLE `articles`
  MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT для таблицы `categories`
--
ALTER TABLE `categories`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT для таблицы `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `articles`
--
ALTER TABLE `articles`
  ADD CONSTRAINT `articles_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
