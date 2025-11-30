import { useQuery } from "@tanstack/react-query";
import normalBulb from "/assets/bulbs/normalBulb.png"
import bulkyBulb from "/assets/bulbs/bulkyBulb.png"
import vintageBulb from "/assets/bulbs/vintageBulb.png"
import slimBulb from "/assets/bulbs/slimBulb.png"
import gu10Bulb from "/assets/bulbs/gu10Bulb.png"

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

export default function Bulbs() {
    const { data: bulbs, isLoading, error } = useQuery({
        queryKey: ['bulbs'],
        queryFn: fetchAllBulbs,
    });

    console.log(bulbs)
    return (
        <div className="bulbs-page">
            {bulbs.map(bulb =>
                <section className="bulb-item" key={bulb.ip}>
                    <img src={getBulbImg(bulb.type)} alt={`${bulb.type} wiz light bulb`} />
                    <section className="bulb-item-info">
                        <h2 className="bulb-name">
                            {bulb.name}
                        </h2>
                        <h2 className="bulb-mac">
                            {bulb.mac.toUpperCase()}
                        </h2>
                        <button className="bulb-flash-button">
                            Flash Light
                        </button>
                    </section>
                </section>
            )}
        </div>
    )
}
