import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

const fetchRules = async () => {
    const response = await fetch("http://localhost:10100/api/rules");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

export default function Rules() {
    const navigate = useNavigate()

    const { data: rules, isLoading: loadingRules, error: errorRules } = useQuery({
        queryKey: ['rules'],
        queryFn: fetchRules,
    });

    return (
        <div className="rules-page">
            <section className="current-rules">
                {rules ? rules : "Found no rules. Click the button to add some."}
            </section>
            <button
                className="add-rule-button"
                onClick={() => navigate("/AddRules")}
            >
                Add a Rule
            </button>
        </div>
    )
}
