import { useMutation, useQuery } from "@tanstack/react-query";
import normalBulb from "/assets/bulbs/normalBulb.png"
import bulkyBulb from "/assets/bulbs/bulkyBulb.png"
import vintageBulb from "/assets/bulbs/vintageBulb.png"
import slimBulb from "/assets/bulbs/slimBulb.png"
import gu10Bulb from "/assets/bulbs/gu10Bulb.png"
import { useState } from "react";

const fetchAllBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs");
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json();
}

const fetchNewBulbs = async () => {
    const response = await fetch("http://localhost:10100/api/bulbs/refresh");
    if (!response.ok) throw new Error('Refresh bulbs response was not ok');
    return response.json();
}

const flashBulbRequest = async (mac) => {
    const response = await fetch("http://localhost:10100/api/bulbs/flash", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ mac }),
    });
    if (!response.ok) throw new Error('Flash bulb request failed');
};

const getBulbImg = (bulbType) => {
    const bulbTypeLower = bulbType.toLowerCase();
    if (bulbTypeLower === "normal") return normalBulb;
    if (bulbTypeLower === "bulky") return bulkyBulb;
    if (bulbTypeLower === "slim") return slimBulb;
    if (bulbTypeLower === "gu10") return gu10Bulb;
    if (bulbTypeLower === "vintage") return vintageBulb;
    return normalBulb;
}

export default function Bulbs() {
    const [flashingMac, setFlashingMac] = useState(null);
    //const [changingNameMac, setChangingNameMac] = useState(null)


    const { data: bulbs, error } = useQuery({
        queryKey: ['bulbs'],
        queryFn: fetchAllBulbs,
    });

    const { data: numOfNewBulbs, refetch: refreshBulbs, isFetching } = useQuery({
        queryKey: ['numOfNewBulbs'],
        queryFn: fetchNewBulbs,
        enabled: false,
    });

    const flashBulbMutation = useMutation({
        mutationFn: flashBulbRequest,
    });

    const flashBulb = (mac) => {
        setFlashingMac(mac);
        flashBulbMutation.mutate(mac);

        setTimeout(() => {
            setFlashingMac(null);
        }, 2900);
    }

    return (
        <div className="bulbs-page-contents">
            {isFetching ?
                <h1>
                    Looking for additional bulbs on the network...
                </h1> :
                numOfNewBulbs &&
                <h1>
                    {numOfNewBulbs.num === 0 ? "Found no new light bulbs on the network" : `Found and added ${numOfNewBulbs.num} new bulbs!`}
                </h1>
            }
            <div className="bulbs-page">
                {error ? (
                    <h2>Something went wrong getting bulbs</h2>
                ) :
                    bulbs?.length && bulbs?.map(bulb =>
                        <section className="bulb-item" key={bulb.ip}>
                            <img src={getBulbImg(bulb.type)} alt={`${bulb.type} wiz light bulb`} />
                            <section className="bulb-item-info">
                                <h2 className="bulb-name">
                                    {bulb.name}
                                </h2>
                                <h2 className="bulb-mac">
                                    {bulb.mac.toUpperCase()}
                                </h2>
                                <button
                                    className="bulb-flash-button"
                                    onClick={() => flashBulb(bulb.mac)}
                                    disabled={flashingMac === bulb.mac}
                                >
                                    {flashingMac === bulb.mac ? "Flashing..." : "Flash Light"}
                                </button>
                            </section>
                        </section>
                    )}
            </div>
            <button
                className="bulbs-refresh-button"
                onClick={() => refreshBulbs()}
                disabled={isFetching}
            >
                {isFetching ? "Searching..." : "Add Bulbs"}
            </button>
        </div>
    )
}
