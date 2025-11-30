import { NavLink, Outlet } from "react-router-dom";

export default function Navbar() {
    return (
        <>
            <header className="nav-header">
                <div className="nav-links">
                    <NavLink to="plexinfo">
                        Plex
                    </NavLink>
                    <NavLink to="rules">
                        Rules
                    </NavLink>
                    <NavLink to="bulbs">
                        Bulbs
                    </NavLink>
                </div>
            </header>
            <div className="contents">
                <Outlet />
            </div>
        </>
    )
}
