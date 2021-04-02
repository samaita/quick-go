
import React, { useEffect, useState } from 'react'
import firebase from 'firebase/app'
import 'firebase/auth'

import './App.css'

function App() {

	useEffect(() => {
		let firebaseConfig = {
			apiKey: process.env.REACT_APP_QUICK_GO_FIREBASE_API_KEY,
			authDomain: process.env.REACT_APP_QUICK_GO_FIREBASE_AUTH_DOMAIN,
			projectId: process.env.REACT_APP_QUICK_GO_FIREBASE_PROJECT_ID,
			storageBucket: process.env.REACT_APP_QUICK_GO_FIREBASE_STORAGE_BUCKET,
			messagingSenderId: process.env.REACT_APP_QUICK_GO_FIREBASE_MESSAGING_SENDER_ID,
			appId: process.env.REACT_APP_QUICK_GO_FIREBASE_APP_ID
		};
		firebase.initializeApp(firebaseConfig);
	})

	const handleLoginWithGoggle = () => {
		var provider = new firebase.auth.GoogleAuthProvider();
		provider.addScope("profile");
		provider.addScope("email");
		provider.addScope("https://www.googleapis.com/auth/plus.me");
		firebase.auth().signInWithPopup(provider)
			.then((result) => {
				/** @type {firebase.auth.OAuthCredential} */
				var credential = result.credential;

				// This gives you a Google Access Token. You can use it to access the Google API.
				var token = credential.accessToken;
				// The signed-in user info.
				var user = result.user;

				console.log(user, token, credential)
				// ...
			}).catch((error) => {
				// Handle Errors here.
				var errorCode = error.code;
				var errorMessage = error.message;
				// The email of the user's account used.
				var email = error.email;
				// The firebase.auth.AuthCredential type that was used.
				var credential = error.credential;
				// ...

				console.log(error)
			});
	}

	return (
		<div>
			<LoginPage
				handleLoginWithGoogle={handleLoginWithGoggle} />
		</div>
	)
}

const LoginPage = ({
	handleLoginWithGoogle
}) => {
	return (
		<div className="bg-gray-200">
			<div className="flex h-screen justify-center items-center">
				<div className="m-auto w-5/12 h-96 flex shadow-md rounded-md">
					<div className="w-5/12 h-full bg-gray-100 rounded-l-md p-5">
						<span className="w-24 h-24 block m-auto mb-3 text-green-500">
							<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
								<path d="M10 3.5a1.5 1.5 0 013 0V4a1 1 0 001 1h3a1 1 0 011 1v3a1 1 0 01-1 1h-.5a1.5 1.5 0 000 3h.5a1 1 0 011 1v3a1 1 0 01-1 1h-3a1 1 0 01-1-1v-.5a1.5 1.5 0 00-3 0v.5a1 1 0 01-1 1H6a1 1 0 01-1-1v-3a1 1 0 00-1-1h-.5a1.5 1.5 0 010-3H4a1 1 0 001-1V6a1 1 0 011-1h3a1 1 0 001-1v-.5z" />
							</svg>
						</span>
						<form>
							<input className="w-full h-9 mb-2 p-2 shadow-inner rounded-md text-gray-600 font-light focus:border-transparent" type="email" required pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$" placeholder="your@email.com" />
							<input className="w-full h-9 mb-2 p-2 shadow-inner rounded-md text-gray-600 font-light focus:border-green-500" type="password" required placeholder="password" />
							<button className="bg-green-500 text-white text-md font-light w-full p-2 rounded-md mb-1" type="submit">Login</button>
						</form>
						<span className="text-gray-600 font-light text-xs text-center w-full block mb-2">----- or login with -----</span>
						<button className="bg-red-500 text-white text-sm font-light w-full p-2 rounded-md mb-2" onClick={() => handleLoginWithGoogle()}>Google</button>
						<button className="bg-blue-500 text-white text-sm font-light w-full p-2 rounded-md">Facebook</button>
					</div>
					<div className="w-7/12 h-full bg-green-500 rounded-r-md">

					</div>
				</div>
			</div>
		</div>
	)
}

export default App
