import { useQuery } from "@tanstack/react-query";

import Accounts from "./Accounts";

import normalBulb from "/assets/bulbs/normalBulb.png"
import bulkyBulb from "/assets/bulbs/bulkyBulb.png"
import vintageBulb from "/assets/bulbs/vintageBulb.png"
import slimBulb from "/assets/bulbs/slimBulb.png"
import gu10Bulb from "/assets/bulbs/gu10Bulb.png"
import Devices from "./Devices";

const fetchRules = async () => {
    const response = await fetch("http://localhost:10100/api/rules");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchAccounts = async () => {
    const response = await fetch("http://localhost:10100/api/accounts");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchDevices = async () => {
    const response = await fetch("http://localhost:10100/api/devices");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
};

const fetchAllBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
}

const getBulbImg = (bulbType) => {
    const bulbTypeLower = bulbType.toLowerCase();
    if (bulbTypeLower === "normal") return normalBulb;
    if (bulbTypeLower === "bulky") return bulkyBulb;
    if (bulbTypeLower === "slim") return slimBulb;
    if (bulbTypeLower === "gu10") return gu10Bulb;
    if (bulbTypeLower === "vintage") return vintageBulb;
    return normalBulb;
}

export default function Rules() {
    const { data: rules, isLoading: loadingRules, error: errorRules } = useQuery({
        queryKey: ['rules'],
        queryFn: fetchRules,
    });

    const { data: bulbs, isLoading: loadingBulbs, error: errorBulbs } = useQuery({
        queryKey: ['bulbs'],
        queryFn: fetchAllBulbs,
    });

    const { data: savedAccounts, isLoading: loadingSavedAccounts, error: errorSavedAccounts } = useQuery({
        queryKey: ['accounts'],
        queryFn: fetchAccounts,
    });

    const { data: savedDevices, isLoading: loadingSavedDevices, error: errorSavedDevices } = useQuery({
        queryKey: ['devices'],
        queryFn: fetchDevices,
    });
    console.log(rules)

    return (
        <div className="rules-page">
            <Accounts data={savedAccounts} />
            <Devices data={savedDevices} />
        </div>
    )
}
