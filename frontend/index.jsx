import { createRoot } from "react-dom/client"
import { BrowserRouter, Routes, Route } from "react-router-dom"

import Navbar from "./src/pages/Navbar"

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navbar />}/>
            </Routes>
        </BrowserRouter>
    )
}

createRoot(document.getElementById("root")).render(
    <App />
)
