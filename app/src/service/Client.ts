import Data from "../model/Data";

class Client {
    private static _latest: Data;
    private static _median: Data;
    private static _max: Data;
    private static _min: Data;
    private static _today: Array<Data>;
    private static type: string;
    private static client: Client = new Client();

    private constructor() {
        const dt = Data.getRandomDataType();

        Client._latest = Data.build(dt);
        Client._median = Data.build(dt);
        Client._max = Data.build(dt);
        Client._min = Data.build(dt);

        // for (let i = 0; i < 3; i++) {
        //     Client._today.push(Data.build());
        // }
    }

    public static getInstance() {
        return this.client;
    }

    public get latest() {
        return Client._latest;
    }

    public get median() {
        return Client._median;
    }

    public get max() {
        return Client._max;
    }

    public get min() {
        return Client._min;
    }

    public get today() {
        return Client._today;
    }
}

export default Client;