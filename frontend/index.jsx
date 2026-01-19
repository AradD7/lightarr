import { createRoot } from "react-dom/client"
import { BrowserRouter, Routes, Route } from "react-router-dom"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"

import Navbar from "./src/pages/Navbar"
import Rules from "./src/pages/Rules/Rules"
import PlexPage from "./src/pages/PlexPage"
import Bulbs from "./src/pages/Bulbs"
import Welcome from "./src/pages/Welcome"
import AddRule from "./src/pages/Rules/AddRule"

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navbar />}>
                    <Route index element={<Welcome />} />

                    <Route path="Rules" element={<Rules />} />
                    <Route path="AddRules" element={<AddRule />} />
                    <Route path="PlexInfo" element={<PlexPage />} />
                    <Route path="Bulbs" element={<Bulbs />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}

const queryClinet = new QueryClient();

createRoot(document.getElementById("root")).render(
    <QueryClientProvider client={queryClinet}>
        <App />
    </QueryClientProvider>
)
