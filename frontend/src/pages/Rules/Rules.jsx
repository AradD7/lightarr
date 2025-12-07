import { useQuery } from "@tanstack/react-query";

import AddRule from "./AddRule";


const fetchRules = async () => {
    const response = await fetch("http://localhost:10100/api/rules");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

export default function Rules() {
    const { data: rules, isLoading: loadingRules, error: errorRules } = useQuery({
        queryKey: ['rules'],
        queryFn: fetchRules,
    });

    console.log(rules)

    return (
        <div className="rules-page">
            <AddRule />
        </div>
    )
}
