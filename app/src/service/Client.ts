import {Data, DataType} from "../model/Data";

class Client {
    private static _latest: Data;
    private static _median: Data;
    private static _max: Data;
    private static _min: Data;
    private static _today: Array<Data> = new Array<Data>();
    private static type: DataType = DataType.CARBON_MONOXIDE;
    private static client: Client = new Client();

    private constructor() {
        Client._latest = Data.build(Client.type);
        Client._median = Data.build(Client.type);
        Client._max = Data.build(Client.type);
        Client._min = Data.build(Client.type);

        for (let i = 0; i < 10; i++) {
            Client._today.push(Data.build(Client.type));
        }

        Client._today.sort((a: Data, b: Data) => {
            return +b.timestamp - +a.timestamp;
        });
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