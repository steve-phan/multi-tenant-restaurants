"use client";

import { useState } from 'react';
import api from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';

interface Table {
  ID: string;
  Name: string;
  Capacity: number;
}

export default function TablesPage() {
  const [tables, setTables] = useState<Table[]>([]);
  const [restaurantId, setRestaurantId] = useState<string | null>(null);

  const fetchTables = async () => {
    if (!restaurantId) return;
    try {
        const res = await api.get(`/public/restaurants/${restaurantId}/tables/available`);
        setTables(res.data);
    } catch (e) {
        console.error(e);
    }
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Table Management</h1>
      
      <div className="p-4 border rounded bg-yellow-50 mb-6">
        <Input 
            placeholder="Enter Restaurant ID" 
            value={restaurantId || ''} 
            onChange={(e) => setRestaurantId(e.target.value)} 
            className="max-w-sm mt-2"
        />
        <Button onClick={fetchTables} className="mt-2">Load Tables</Button>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle>Tables</CardTitle>
            <AddTableDialog restaurantId={restaurantId} onSuccess={fetchTables} />
        </CardHeader>
        <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {tables.map(t => (
                    <div key={t.ID} className="p-4 border rounded flex flex-col items-center justify-center bg-gray-50">
                        <span className="font-bold text-lg">{t.Name}</span>
                        <span className="text-gray-500">{t.Capacity} Seats</span>
                    </div>
                ))}
            </div>
        </CardContent>
      </Card>
    </div>
  );
}

function AddTableDialog({ restaurantId, onSuccess }: { restaurantId: string | null, onSuccess: () => void }) {
    const [name, setName] = useState('');
    const [capacity, setCapacity] = useState('');
    const [open, setOpen] = useState(false);

    const handleSubmit = async () => {
        if (!restaurantId) return;
        await api.post(`/restaurants/${restaurantId}/tables`, { 
            Name: name, 
            Capacity: parseInt(capacity) 
        });
        setOpen(false);
        setName('');
        setCapacity('');
        onSuccess();
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild><Button size="sm">Add New</Button></DialogTrigger>
            <DialogContent>
                <DialogHeader><DialogTitle>Add Table</DialogTitle></DialogHeader>
                <div className="space-y-4">
                    <div className="space-y-2">
                        <Label>Name</Label>
                        <Input value={name} onChange={e => setName(e.target.value)} />
                    </div>
                    <div className="space-y-2">
                        <Label>Capacity</Label>
                        <Input type="number" value={capacity} onChange={e => setCapacity(e.target.value)} />
                    </div>
                    <Button onClick={handleSubmit}>Save</Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}
