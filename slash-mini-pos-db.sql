-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               8.0.30 - MySQL Community Server - GPL
-- Server OS:                    Win64
-- HeidiSQL Version:             12.1.0.6537
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for slash-pos
DROP DATABASE IF EXISTS `slash-pos`;
CREATE DATABASE IF NOT EXISTS `slash-pos` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `slash-pos`;

-- Dumping structure for table slash-pos.orders
DROP TABLE IF EXISTS `orders`;
CREATE TABLE IF NOT EXISTS `orders` (
  `id` varchar(50) NOT NULL,
  `user_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `customer_name` varchar(255) NOT NULL,
  `customer_phone` varchar(20) DEFAULT NULL,
  `customer_address` text,
  `total_amount` varchar(50) NOT NULL,
  `order_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `expired` datetime NOT NULL,
  `status` enum('pending','done') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `FK_orders_users` (`user_id`),
  CONSTRAINT `FK_orders_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table slash-pos.orders: ~0 rows (approximately)

-- Dumping structure for table slash-pos.order_items
DROP TABLE IF EXISTS `order_items`;
CREATE TABLE IF NOT EXISTS `order_items` (
  `order_item_id` varchar(50) NOT NULL,
  `order_id` varchar(50) NOT NULL,
  `product_id` varchar(50) NOT NULL,
  `quantity` int NOT NULL,
  `price` varchar(50) NOT NULL DEFAULT '0',
  PRIMARY KEY (`order_item_id`),
  KEY `order_id` (`order_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table slash-pos.order_items: ~2 rows (approximately)

-- Dumping structure for table slash-pos.products
DROP TABLE IF EXISTS `products`;
CREATE TABLE IF NOT EXISTS `products` (
  `id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `color` varchar(255) DEFAULT NULL,
  `size` varchar(255) DEFAULT NULL,
  `stock` int DEFAULT NULL,
  `price` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `color` (`color`),
  KEY `size` (`size`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table slash-pos.products: ~11 rows (approximately)
INSERT INTO `products` (`id`, `name`, `color`, `size`, `stock`, `price`) VALUES
	('item-A1b2C3d4E5f6G7h', 'Office Suit', 'Black', 'XL', 32, 'Rp 1.750.000'),
	('item-K7j6L5k4H3n2J1m', 'Denim Jacket', 'Dark Blue', 'XL', 30, 'Rp 650.000'),
	('item-L4k3J2h1F9g8H7t', 'Hoodie', 'Grey', 'M', 60, 'Rp 320.000'),
	('item-M3n2B1v9C8x7D6w', 'T-Shirt', 'Black', 'S', 120, 'Rp 150.000'),
	('item-N5m4B3v2C1z9X8w', 'Formal Shirt', 'Light Blue', 'L', 85, 'Rp 400.000'),
	('item-O5p4L3n2K1m9J8h', 'Blazer', 'Navy', 'L', 54, 'Rp 1.200.000'),
	('item-P8o7I6u5Y4t3R2e', 'Jeans', 'Blue', 'XL', 60, 'Rp 450.000'),
	('item-Q8w7E6r5T4y3U2i', 'Short Pants', 'Khaki', 'S', 69, 'Rp 200.000'),
	('item-R9e8T7w6Y5u4I3o', 'Jogger Pants', 'Black', 'M', 70, 'Rp 300.000'),
	('item-X2z1A9v8B7c6N5m', 'Sweater', 'Beige', 'XL', 50, 'Rp 350.000'),
	('item-Z9y8X7w6V5u4T3s', 'Casual Shirt', 'White', 'L', 100, 'Rp 250.000');

-- Dumping structure for table slash-pos.users
DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(50) NOT NULL,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `role` enum('helper') NOT NULL DEFAULT 'helper',
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table slash-pos.users: ~4 rows (approximately)
INSERT INTO `users` (`id`, `username`, `email`, `password`, `role`, `is_active`, `created_at`, `updated_at`) VALUES
	('user-HYEFmLONNJhrj17y', 'testuser2', 'test@user2.com', '$2a$10$MyeAEN5UrLQqa4zPii2b..n/WIdQjMQjPENzYxi1KlQm7mOKizP7S', 'helper', 1, '2025-01-08 10:51:14', '2025-01-08 10:51:14'),
	('user-LVpjmO25MulBGxFI', 'testuser5', 'test@user5.com', '$2a$10$sazz0r4DHGV5R/xWnOi7l.lUiMPNQ0A9HTq4/ZSLiPlmUg6IsNgZm', 'helper', 1, '2025-01-10 03:14:38', '2025-01-10 03:14:38'),
	('user-sqVCEOQcXzwtMuMF', 'testuser', 'test@user.com', '$2a$10$fV5I1aBSbCRMjZCDwQ67j.XWX0lCEQ7C00l2cl2jNACMn2Yr9E16S', 'helper', 1, '2025-01-08 10:50:46', '2025-01-08 10:50:46'),
	('user-yB90wcffptjHauUP', 'testuser3', 'test@user3.com', '$2a$10$SGDnxVvIb.fpo5d/4ez4PewPf09Y0cDz3.wwhGNvOyCTfNzb66rzG', 'helper', 1, '2025-01-09 10:22:12', '2025-01-09 10:22:12');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
