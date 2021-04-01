
import React, { useEffect } from 'react'
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

	const LoginWithGoggle = () => {
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
			<h1>Hello World</h1>
			<a href="#" onClick={() => LoginWithGoggle()}>LOGIN</a>
		</div>
	)
}

export default App
