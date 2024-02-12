package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// Adjust the import path based on your project structure
)

func AdminPageHandler(c *gin.Context) {
	html := `

	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>clodevo Proxy</title>
		<link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap" rel="stylesheet">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.3/css/all.min.css">
		<style>
			body {
				font-family: 'Roboto', sans-serif;
				margin: 0;
				padding: 0;
				background: #fff;
				color: #333;
				overflow-x: hidden;
				display: flex;
				flex-direction: column;
				height: 100vh;
			}
	
			.hero {
				background-color: #fff;
				color: #333;
				padding: 50px 20px;
				text-align: center;
				position: relative;
				overflow: hidden;
			}
	
			/* Particle animation */
			.particle {
				position: absolute;
				border-radius: 70%;
				background: rgba(255, 165, 0, 0.6);
				animation: float 10s linear infinite;
				opacity: 0.6;
			}
	
			@keyframes float {
				0%, 100% { transform: translateY(0) translateX(0); }
				50% { transform: translateY(-20px) translateX(15px); }
			}
	
			/* Generate multiple particles */
			.hero::before, .hero::after {
				content: '';
				position: absolute;
				width: 10px;
				height: 10px;
				background: rgba(255, 165, 0, 0.6);
				border-radius: 50%;
				top: 20%;
				left: 25%;
				animation: float 8s linear infinite;
			}
	
			.hero::after {
				top: 70%;
				left: 80%;
				animation-duration: 12s;
			}
	
			h1, p, .buttons {
				position: relative;
				z-index: 1;
			}
	
			.buttons {
				display: flex;
				justify-content: center;
				gap: 20px;
				margin-top: 20px;
			}
	
			.button {
				background-color: #ffa500;
				color: #fff;
				padding: 15px 30px;
				text-decoration: none;
				border-radius: 5px;
				display: inline-block;
				font-size: 18px;
				font-weight: bold;
				transition: background-color 0.3s ease, transform 0.2s ease;
			}
	
			.button:hover {
				background-color: #ff8c00;
				transform: scale(1.05);
			}
	
			.container {
				padding: 20px;
				display: flex;
				flex-wrap: wrap;
				justify-content: space-around;
				align-content: flex-start;
				overflow: auto;
				flex-grow: 1;
			}
	
			.feature {
				background: #f9f9f9;
				box-shadow: 0 2px 15px rgba(0,0,0,0.1);
				border-radius: 8px;
				padding: 10px;
				margin: 5px;
				flex: 1 1 calc(30% - 20px);
				display: flex;
				align-items: center;
				justify-content: center;
				flex-direction: column;
				min-width: 140px;
				max-height: 160px;
			}
	
			.feature-icon {
				font-size: 20px;
				color: #ffa500;
				margin-bottom: 5px;
			}
	
			.footer {
				text-align: center;
				padding: 20px;
				margin-top: auto;
				background: #f1f1f1;
			}
	
			.footer a {
				color: #ffa500;
				text-decoration: none;
			}
		</style>
	</head>
	<body>
		<div class="hero">
			<div class="particle" style="width: 15px; height: 15px; top: 10%; left: 50%;"></div>
			<div class="particle" style="width: 10px; height: 10px; top: 50%; left: 75%; animation-duration: 15s;"></div>
			<div class="particle" style="width: 5px; height: 5px; top: 75%; left: 20%; animation-duration: 7s;"></div>
			<h1>Welcome to clodevo Proxy</h1>
			<p>Your cutting-edge forward proxy solution.</p>
			<div class="buttons">
				<a href="/swagger/index.html" class="button">Experience the Power</a>
				<a href="https://www.clodevo.com" class="button">Contact</a>
				<a href="https://github.com/clodevo/clodevo-http-proxy" class="button github"><i class="fab fa-github fa-2x"></i></a>
			</div>
		</div>
		<div class="container">
			<div class="feature">
				<i class="feature-icon fas fa-forward"></i>
				<h2>Forward Proxy</h2>
				<p>Mediates traffic between users and the internet.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-users"></i>
				<h2>Multi-Tenant Support</h2>
				<p>Isolates data and configurations for each tenant.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-user-shield"></i>
				<h2>Tenant Authentication</h2>
				<p>Ensures secure access with robust authentication.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-eye"></i>
				<h2>Transparent Proxy</h2>
				<p>Seamless integration without client-side config.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-chart-line"></i>
				<h2>Logging and Monitoring</h2>
				<p>Comprehensive insights into network activities.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-tachometer-alt"></i>
				<h2>Scalability and Performance</h2>
				<p>Optimized for high load conditions.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-lock"></i>
				<h2>Security Enhancements</h2>
				<p>Advanced features like SSL/TLS encryption.</p>
			</div>
			<div class="feature">
				<i class="feature-icon fas fa-cloud"></i>
				<h2>Kubernetes Native Support</h2>
				<p>Seamless integration with Kubernetes environments.</p>
			</div>
		</div>
	</body>
	</html>	
    `
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
