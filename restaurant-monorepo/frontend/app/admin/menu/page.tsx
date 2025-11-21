"use client";

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';

// Types
interface Category {
  ID: string;
  Name: string;
}

interface MenuItem {
  ID: string;
  Name: string;
  Price: number;
  CategoryID: string;
}

export default function MenuPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [items, setItems] = useState<MenuItem[]>([]);
  const [restaurantId, setRestaurantId] = useState<string | null>(null); // In real app, fetch from user context/api

  // Temporary: Create a restaurant if none exists (for MVP flow)
  useEffect(() => {
    const fetchRestaurant = async () => {
       // This is a hack for MVP to get a restaurant ID. 
       // In a real app, we'd get the user's restaurant.
       // For now, we'll just create one if we don't have one stored or fetch the first one.
       // Let's assume the user has one.
       // TODO: Implement proper restaurant fetching
    };
    fetchRestaurant();
  }, []);

  const fetchMenu = async () => {
    if (!restaurantId) return;
    // Fetch categories and items
    // This requires endpoints that return lists. 
    // Our backend currently has public endpoints for this.
    try {
        const res = await api.get(`/public/restaurants/${restaurantId}/menu`);
        setCategories(res.data.categories);
        setItems(res.data.items);
    } catch (e) {
        console.error(e);
    }
  };

  // ... (Implementation of Add Category and Add Item forms)
  // Since I need to know the restaurant ID, I'll add a simple input for it for now or auto-create one.
  
  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Menu Management</h1>
      <p className="mb-4 text-gray-500">Manage your menu categories and items.</p>
      
      {/* Placeholder for MVP */}
      <div className="p-4 border rounded bg-yellow-50 mb-6">
        <p>Note: For this MVP, please ensure you have created a restaurant via API or use the ID provided.</p>
        <Input 
            placeholder="Enter Restaurant ID to manage" 
            value={restaurantId || ''} 
            onChange={(e) => setRestaurantId(e.target.value)} 
            className="max-w-sm mt-2"
        />
        <Button onClick={fetchMenu} className="mt-2">Load Menu</Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Categories */}
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle>Categories</CardTitle>
                <AddCategoryDialog restaurantId={restaurantId} onSuccess={fetchMenu} />
            </CardHeader>
            <CardContent>
                <ul>
                    {categories.map(c => (
                        <li key={c.ID} className="p-2 border-b">{c.Name}</li>
                    ))}
                </ul>
            </CardContent>
        </Card>

        {/* Items */}
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle>Items</CardTitle>
                <AddItemDialog restaurantId={restaurantId} categories={categories} onSuccess={fetchMenu} />
            </CardHeader>
            <CardContent>
                <ul>
                    {items.map(i => (
                        <li key={i.ID} className="p-2 border-b flex justify-between">
                            <span>{i.Name}</span>
                            <span>${i.Price}</span>
                        </li>
                    ))}
                </ul>
            </CardContent>
        </Card>
      </div>
    </div>
  );
}

function AddCategoryDialog({ restaurantId, onSuccess }: { restaurantId: string | null, onSuccess: () => void }) {
    const [name, setName] = useState('');
    const [open, setOpen] = useState(false);

    const handleSubmit = async () => {
        if (!restaurantId) return;
        await api.post(`/restaurants/${restaurantId}/menu/categories`, { Name: name });
        setOpen(false);
        setName('');
        onSuccess();
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild><Button size="sm">Add New</Button></DialogTrigger>
            <DialogContent>
                <DialogHeader><DialogTitle>Add Category</DialogTitle></DialogHeader>
                <div className="space-y-4">
                    <div className="space-y-2">
                        <Label>Name</Label>
                        <Input value={name} onChange={e => setName(e.target.value)} />
                    </div>
                    <Button onClick={handleSubmit}>Save</Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}

function AddItemDialog({ restaurantId, categories, onSuccess }: { restaurantId: string | null, categories: Category[], onSuccess: () => void }) {
    const [name, setName] = useState('');
    const [price, setPrice] = useState('');
    const [catId, setCatId] = useState('');
    const [open, setOpen] = useState(false);

    const handleSubmit = async () => {
        if (!restaurantId) return;
        await api.post(`/restaurants/${restaurantId}/menu/items`, { 
            Name: name, 
            Price: parseFloat(price), 
            CategoryID: catId 
        });
        setOpen(false);
        setName('');
        setPrice('');
        onSuccess();
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild><Button size="sm">Add New</Button></DialogTrigger>
            <DialogContent>
                <DialogHeader><DialogTitle>Add Item</DialogTitle></DialogHeader>
                <div className="space-y-4">
                    <div className="space-y-2">
                        <Label>Name</Label>
                        <Input value={name} onChange={e => setName(e.target.value)} />
                    </div>
                    <div className="space-y-2">
                        <Label>Price</Label>
                        <Input type="number" value={price} onChange={e => setPrice(e.target.value)} />
                    </div>
                    <div className="space-y-2">
                        <Label>Category</Label>
                        <select className="w-full border p-2 rounded" value={catId} onChange={e => setCatId(e.target.value)}>
                            <option value="">Select Category</option>
                            {categories.map(c => <option key={c.ID} value={c.ID}>{c.Name}</option>)}
                        </select>
                    </div>
                    <Button onClick={handleSubmit}>Save</Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}
